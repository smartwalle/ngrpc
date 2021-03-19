package etcd

import (
	"github.com/smartwalle/etcd4go"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
)

const (
	kDefaultScheme = "etcd"
)

type Registry struct {
	scheme     string
	client     *etcd4go.Client
	watcher    *etcd4go.Watcher
	clientConn resolver.ClientConn
}

func NewRegistry(client *clientv3.Client) *Registry {
	return NewRegistryWithScheme(kDefaultScheme, client)
}

func NewRegistryWithScheme(scheme string, client *clientv3.Client) *Registry {
	var nRegistry = &Registry{scheme: scheme, client: etcd4go.NewClient(client)}
	resolver.Register(nRegistry)
	return nRegistry
}

func (this *Registry) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	this.clientConn = cc
	var key = target.Scheme + "://" + filepath.Join(target.Authority, target.Endpoint)
	this.watcher = this.client.Watch(key, this.watch, clientv3.WithPrefix())
	return this, nil
}

func (this *Registry) watch(watcher *etcd4go.Watcher, event, key, path string, value []byte) {
	var paths = watcher.Values()
	var addrList = make([]resolver.Address, 0, len(paths))
	for _, nValue := range paths {
		var addr = resolver.Address{Addr: string(nValue)}
		addrList = append(addrList, addr)
	}
	this.clientConn.UpdateState(resolver.State{Addresses: addrList})
}

func (this *Registry) Scheme() string {
	return this.scheme
}

func (this *Registry) ResolveNow(option resolver.ResolveNowOptions) {
}

func (this *Registry) Close() {
	if this.watcher != nil {
		this.watcher.Close()
	}
}

func (this *Registry) Register(domain, service, node, addr string, ttl int64) (key string, err error) {
	_, key, err = this.client.Register(this.BuildPath(domain, service, node), addr, ttl)
	return key, err
}

func (this *Registry) Deregister(domain, service, node string) (err error) {
	return this.client.Deregister(this.BuildPath(domain, service, node))
}

func (this *Registry) BuildPath(domain, service, node string) string {
	var target = this.scheme + "://" + filepath.Join(domain, service, node)
	return target
}
