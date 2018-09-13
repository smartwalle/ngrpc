package grpc4go

import (
	"github.com/smartwalle/etcd4go"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
)

const (
	k_DEFAULT_SCHEME = "grpc"
)

type Resolver struct {
	scheme string
	c      *etcd4go.Client
}

func NewResolver(etcd *etcd4go.Client) *Resolver {
	return NewResolverWithScheme(k_DEFAULT_SCHEME, etcd)
}

func NewResolverWithScheme(scheme string, c *etcd4go.Client) *Resolver {
	return &Resolver{scheme: scheme, c: c}
}

func (this *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var key = target.Scheme + "://" + filepath.Join(target.Authority, target.Endpoint)
	watchInfo := this.c.Watch(key, clientv3.WithPrefix())

	watchInfo.Handle(func(eventType, key, path string, value []byte) {
		var paths = watchInfo.GetPaths()
		var addList = make([]resolver.Address, 0, len(paths))
		for _, value := range paths {
			var add = resolver.Address{Addr: string(value)}
			addList = append(addList, add)
		}
		cc.NewAddress(addList)
	})
	return this, nil
}

func (this *Resolver) Scheme() string {
	return this.scheme
}

func (this *Resolver) ResolveNow(option resolver.ResolveNowOption) {
}

func (this *Resolver) Close() {
}

func (this *Resolver) Register(service, node, addr string, ttl int64) (key string, err error) {
	return this.c.RegisterWithKey(this.scheme+"://"+filepath.Join(service, node), addr, ttl)
}

func (this *Resolver) UnRegister(service, node, addr string) (err error) {
	return this.c.RevokeWithKey(this.scheme + "://" + filepath.Join(service, node))
}
