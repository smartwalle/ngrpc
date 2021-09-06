package etcd

import (
	"context"
	"github.com/smartwalle/etcd4go"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
	"sync"
)

const (
	kDefaultScheme = "etcd"
)

type Resolver struct {
	scheme   string
	client   *etcd4go.Client
	mu       *sync.Mutex
	watchers map[string]*etcd4go.Watcher
}

func NewResolver(client *clientv3.Client) *Resolver {
	return NewResolverWithScheme(kDefaultScheme, client)
}

func NewResolverWithScheme(scheme string, client *clientv3.Client) *Resolver {
	var nResolver = &Resolver{scheme: scheme, client: etcd4go.NewClient(client)}
	nResolver.mu = &sync.Mutex{}
	nResolver.watchers = make(map[string]*etcd4go.Watcher)
	resolver.Register(nResolver)
	return nResolver
}

func (this *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var key = target.Scheme + "://" + filepath.Join(target.Authority, target.Endpoint)
	var watcher = this.client.Watch(context.Background(), key, this.watch(cc), clientv3.WithPrefix())
	this.mu.Lock()
	this.watchers[key] = watcher
	this.mu.Unlock()
	return this, nil
}

func (this *Resolver) watch(cc resolver.ClientConn) func(watcher *etcd4go.Watcher, event, key, path string, value []byte) {
	return func(watcher *etcd4go.Watcher, event, key, path string, value []byte) {
		var paths = watcher.Values()
		var addrList = make([]resolver.Address, 0, len(paths))
		for _, nValue := range paths {
			var addr = resolver.Address{Addr: string(nValue)}
			addrList = append(addrList, addr)
		}
		cc.UpdateState(resolver.State{Addresses: addrList})
	}
}

func (this *Resolver) Scheme() string {
	return this.scheme
}

func (this *Resolver) ResolveNow(option resolver.ResolveNowOptions) {
}

func (this *Resolver) Close() {
	this.mu.Lock()
	for _, watcher := range this.watchers {
		if watcher != nil {
			watcher.Close()
		}
	}
	this.mu.Unlock()
}

func (this *Resolver) Register(ctx context.Context, domain, service, node, addr string, ttl int64) (key string, err error) {
	_, key, err = this.client.Register(ctx, this.BuildPath(domain, service, node), addr, ttl)
	return key, err
}

func (this *Resolver) Deregister(ctx context.Context, domain, service, node string) (err error) {
	return this.client.Deregister(ctx, this.BuildPath(domain, service, node))
}

func (this *Resolver) BuildPath(domain, service, node string) string {
	var target = this.scheme + "://" + filepath.Join(domain, service, node)
	return target
}
