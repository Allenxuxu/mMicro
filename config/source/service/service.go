package service

import (
	"context"

	"github.com/Allenxuxu/mMicro/client"
	"github.com/Allenxuxu/mMicro/config/source"
	proto "github.com/Allenxuxu/mMicro/config/source/service/proto"
	"github.com/Allenxuxu/mMicro/logger"
)

var (
	DefaultName      = "go.micro.config"
	DefaultNamespace = "global"
	DefaultPath      = ""
)

type service struct {
	serviceName string
	namespace   string
	path        string
	opts        source.Options
	client      proto.ConfigService
}

func (m *service) Read() (set *source.ChangeSet, err error) {
	client := proto.NewConfigService(m.serviceName, client.DefaultClient)
	req, err := client.Read(context.Background(), &proto.ReadRequest{
		Namespace: m.namespace,
		Path:      m.path,
	})
	if err != nil {
		return nil, err
	}

	return toChangeSet(req.Change.ChangeSet), nil
}

func (m *service) Watch() (w source.Watcher, err error) {
	client := proto.NewConfigService(m.serviceName, client.DefaultClient)
	stream, err := client.Watch(context.Background(), &proto.WatchRequest{
		Namespace: m.namespace,
		Path:      m.path,
	})
	if err != nil {
		if logger.V(logger.ErrorLevel, logger.DefaultLogger) {
			logger.Error("watch err: ", err)
		}
		return
	}
	return newWatcher(stream)
}

// Write is unsupported
func (m *service) Write(cs *source.ChangeSet) error {
	return nil
}

func (m *service) String() string {
	return "service"
}

func NewSource(opts ...source.Option) source.Source {
	var options source.Options
	for _, o := range opts {
		o(&options)
	}

	addr := DefaultName
	namespace := DefaultNamespace
	path := DefaultPath

	if options.Context != nil {
		a, ok := options.Context.Value(serviceNameKey{}).(string)
		if ok {
			addr = a
		}

		k, ok := options.Context.Value(namespaceKey{}).(string)
		if ok {
			namespace = k
		}

		p, ok := options.Context.Value(pathKey{}).(string)
		if ok {
			path = p
		}
	}

	s := &service{
		serviceName: addr,
		opts:        options,
		namespace:   namespace,
		path:        path,
	}

	return s
}