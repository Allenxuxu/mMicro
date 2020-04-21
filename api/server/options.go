package server

import (
	"crypto/tls"
	"net/http"

	"github.com/Allenxuxu/mMicro/api/resolver"
)

type Option func(o *Options)

type Options struct {
	EnableCORS bool
	EnableTLS  bool
	ACMEHosts  []string
	TLSConfig  *tls.Config
	Resolver   resolver.Resolver
	Wrappers   []Wrapper
}

type Wrapper func(h http.Handler) http.Handler

func WrapHandler(w Wrapper) Option {
	return func(o *Options) {
		o.Wrappers = append(o.Wrappers, w)
	}
}

func EnableCORS(b bool) Option {
	return func(o *Options) {
		o.EnableCORS = b
	}
}

func ACMEHosts(hosts ...string) Option {
	return func(o *Options) {
		o.ACMEHosts = hosts
	}
}

func EnableTLS(b bool) Option {
	return func(o *Options) {
		o.EnableTLS = b
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

func Resolver(r resolver.Resolver) Option {
	return func(o *Options) {
		o.Resolver = r
	}
}
