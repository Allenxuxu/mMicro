package router

import (
	"github.com/Allenxuxu/mMicro/api/resolver"
	"github.com/Allenxuxu/mMicro/api/resolver/vpath"
	"github.com/Allenxuxu/mMicro/registry"
	"github.com/Allenxuxu/mMicro/registry/mdns"
)

type Options struct {
	Handler  string
	Registry registry.Registry
	Resolver resolver.Resolver
}

type Option func(o *Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		Handler:  "meta",
		Registry: mdns.NewRegistry(),
	}

	for _, o := range opts {
		o(&options)
	}

	if options.Resolver == nil {
		options.Resolver = vpath.NewResolver(
			resolver.WithHandler(options.Handler),
		)
	}

	return options
}

func WithHandler(h string) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

func WithRegistry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func WithResolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}
