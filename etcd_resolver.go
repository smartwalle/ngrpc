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
	scheme    string
	c         *etcd4go.Client
	watchInfo *etcd4go.WatchInfo
}

func NewETCDResolver(c *etcd4go.Client) *ETCDResolver {
	return NewETCDResolverWithScheme(k_DEFAULT_SCHEME, c)
}

func NewETCDResolverWithScheme(scheme string, c *etcd4go.Client) *ETCDResolver {
	return &ETCDResolver{scheme: scheme, c: c}
}

func (this *ETCDResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	var key = target.Scheme + "://" + filepath.Join(target.Authority, target.Endpoint)
	this.watchInfo = this.c.Watch(key, clientv3.WithPrefix())

	this.watchInfo.Handle(func(eventType, key, path string, value []byte) {
		var paths = this.watchInfo.GetPaths()
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
	this.watchInfo.Close()
}

func (this *ETCDResolver) RegisterService(domain, service, node, addr string, ttl int64) (leaseId int64, key string, err error) {
	return this.c.RegisterWithKey(this.scheme+"://"+filepath.Join(domain, service, node), addr, ttl)
}

func (this *ETCDResolver) UnRegisterService(leaseId int64) (err error) {
	return this.c.Revoke(leaseId)
}
