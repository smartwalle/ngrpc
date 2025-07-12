package ngrpc

import "context"

type Registry interface {
	Scheme() string

	Register(ctx context.Context, domain, service, node, addr string, ttl int64) (key string, err error)

	Unregister(ctx context.Context, domain, service, node string) (err error)

	BuildPath(domain, service, node string) string
}
