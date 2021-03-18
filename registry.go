package grpc4go

type Registry interface {
	Scheme() string

	Register(domain, service, node, addr string, ttl int64) (key string, err error)

	Deregister(domain, service, node string) (err error)

	BuildPath(domain, service, node string) string

	Close()
}
