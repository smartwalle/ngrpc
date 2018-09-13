package grpc4go

import (
	"github.com/smartwalle/etcd4go"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
)

const (
	k_SCHEME_ETCD = "etcdv3"
)

type ETCDResolver struct {
	c *etcd4go.Client
}

func NewETCDResolver(c *etcd4go.Client) *ETCDResolver {
	return &ETCDResolver{c: c}
}

func (this *ETCDResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var key = target.Endpoint
	watchInfo := this.c.Watch(key, clientv3.WithPrefix())

	watchInfo.Handle(func(eventType, key, path string, value []byte) {
		var paths = watchInfo.GetPaths()
		var addrList = make([]resolver.Address, 0, len(paths))
		for _, value := range paths {
			var add = resolver.Address{Addr: string(value)}
			addrList = append(addrList, add)
		}
		cc.NewAddress(addrList)
	})
	return this, nil
}

func (this *ETCDResolver) Scheme() string {
	return k_SCHEME_ETCD
}

func (this *ETCDResolver) ResolveNow(option resolver.ResolveNowOption) {
}

func (this *ETCDResolver) Close() {
}

func (this *ETCDResolver) Register() {
	resolver.Register(this)
}

func (this *ETCDResolver) RegisterDefault() {
	resolver.Register(this)
	resolver.SetDefaultScheme(this.Scheme())
}

func (this *ETCDResolver) RegisterService(service, node, addr string, ttl int64) (key string, err error) {
	return this.c.RegisterWithKey(filepath.Join(service, node), addr, ttl)
}

func (this *ETCDResolver) UnRegisterService(service, node, addr string) (err error) {
	return this.c.RevokeWithKey(filepath.Join(service, node))
}
