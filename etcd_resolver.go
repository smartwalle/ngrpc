package grpc4go

import (
	"github.com/smartwalle/etcd4go"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
)

const (
	k_DEFAULT_SCHEME = "etcd"
)

type ETCDResolver struct {
	scheme string
	c      *etcd4go.Client
}

func NewETCDResolver(c *etcd4go.Client) *ETCDResolver {
	return NewETCDResolverWithScheme(k_DEFAULT_SCHEME, c)
}

func NewETCDResolverWithScheme(scheme string, c *etcd4go.Client) *ETCDResolver {
	return &ETCDResolver{scheme: scheme, c: c}
}

func (this *ETCDResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var key = target.Scheme + "://" + filepath.Join(target.Authority, target.Endpoint)
	watchInfo := this.c.Watch(key, clientv3.WithPrefix())

	watchInfo.Handle(func(eventType, key, path string, value []byte) {
		var paths = watchInfo.GetPaths()
		var addrList = make([]resolver.Address, 0, len(paths))
		for _, value := range paths {
			var addr = resolver.Address{Addr: string(value)}
			addrList = append(addrList, addr)
		}
		cc.NewAddress(addrList)
	})
	return this, nil
}

func (this *ETCDResolver) Scheme() string {
	return this.scheme
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
	resolver.SetDefaultScheme(this.scheme)
}

func (this *ETCDResolver) RegisterService(service, node, addr string, ttl int64) (key string, err error) {
	return this.c.RegisterWithKey(this.scheme+"://"+filepath.Join(service, node), addr, ttl)
}

func (this *ETCDResolver) UnRegisterService(service, node, addr string) (err error) {
	return this.c.RevokeWithKey(this.scheme + "://" + filepath.Join(service, node))
}
