package wrapper

import (
	"context"
	"strings"

	"github.com/Allenxuxu/mMicro/client"
	"github.com/Allenxuxu/mMicro/debug/stats"
	"github.com/Allenxuxu/mMicro/debug/trace"
	"github.com/Allenxuxu/mMicro/metadata"
	"github.com/Allenxuxu/mMicro/server"
)

type clientWrapper struct {
	client.Client

	// headers to inject
	headers metadata.Metadata
}

type traceWrapper struct {
	client.Client

	name  string
	trace trace.Tracer
}

var (
	HeaderPrefix = "Micro-"
)

func (c *clientWrapper) setHeaders(ctx context.Context) context.Context {
	// don't overwrite keys
	return metadata.MergeContext(ctx, c.headers, false)
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx = c.setHeaders(ctx)
	return c.Client.Call(ctx, req, rsp, opts...)
}

func (c *clientWrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	ctx = c.setHeaders(ctx)
	return c.Client.Stream(ctx, req, opts...)
}

func (c *clientWrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	ctx = c.setHeaders(ctx)
	return c.Client.Publish(ctx, p, opts...)
}

func (c *traceWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	newCtx, s := c.trace.Start(ctx, req.Service()+"."+req.Endpoint())

	s.Type = trace.SpanTypeRequestOutbound
	err := c.Client.Call(newCtx, req, rsp, opts...)
	if err != nil {
		s.Metadata["error"] = err.Error()
	}

	// finish the trace
	c.trace.Finish(s)

	return err
}

// HandlerStats wraps a server handler to generate request/error stats
func HandlerStats(stats stats.Stats) server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// execute the handler
			err := h(ctx, req, rsp)
			// record the stats
			stats.Record(err)
			// return the error
			return err
		}
	}
}

// TraceCall is a call tracing wrapper
func TraceCall(name string, t trace.Tracer, c client.Client) client.Client {
	return &traceWrapper{
		name:   name,
		trace:  t,
		Client: c,
	}
}

// TraceHandler wraps a server handler to perform tracing
func TraceHandler(t trace.Tracer) server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			// don't store traces for debug
			if strings.HasPrefix(req.Endpoint(), "Debug.") {
				return h(ctx, req, rsp)
			}

			// get the span
			newCtx, s := t.Start(ctx, req.Service()+"."+req.Endpoint())
			s.Type = trace.SpanTypeRequestInbound

			err := h(newCtx, req, rsp)
			if err != nil {
				s.Metadata["error"] = err.Error()
			}

			// finish
			t.Finish(s)

			return err
		}
	}
}
