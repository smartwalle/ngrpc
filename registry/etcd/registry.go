package etcd

import (
	"bytes"
	"context"
	"github.com/smartwalle/etcd4go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
	"sync"
)

const (
	kDefaultScheme = "etcd"
)

type Registry struct {
	scheme   string
	client   *etcd4go.Client
	mu       *sync.Mutex
	watchers map[string]*etcd4go.Watcher
}

func NewRegistry(client *clientv3.Client) *Registry {
	return NewRegistryWithScheme(kDefaultScheme, client)
}

func NewRegistryWithScheme(scheme string, client *clientv3.Client) *Registry {
	var r = &Registry{scheme: scheme, client: etcd4go.NewClient(client)}
	r.mu = &sync.Mutex{}
	r.watchers = make(map[string]*etcd4go.Watcher)
	resolver.Register(r)
	return r
}

func (this *Registry) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	//var key = target.Scheme + "://" + filepath.Join(target.Authority, target.Endpoint)
	var key = this.buildPath(target.Authority, target.Endpoint)
	var watcher = this.client.Watch(context.Background(), key, this.watch(cc), clientv3.WithPrefix())
	this.mu.Lock()
	this.update(cc, watcher.Values())
	this.watchers[key] = watcher
	this.mu.Unlock()
	return this, nil
}

func (this *Registry) watch(cc resolver.ClientConn) func(watcher *etcd4go.Watcher, event, key, path string, value []byte) {
	return func(watcher *etcd4go.Watcher, event, key, path string, value []byte) {
		var values = watcher.Values()
		this.update(cc, values)
	}
}

func (this *Registry) update(cc resolver.ClientConn, values map[string][]byte) {
	var addresses = make([]resolver.Address, 0, len(values))
	for _, value := range values {
		var addr = resolver.Address{Addr: string(value)}
		addresses = append(addresses, addr)
	}
	cc.UpdateState(resolver.State{Addresses: addresses})
}

func (this *Registry) Scheme() string {
	return this.scheme
}

func (this *Registry) ResolveNow(options resolver.ResolveNowOptions) {
}

func (this *Registry) Close() {
	this.mu.Lock()
	for _, watcher := range this.watchers {
		if watcher != nil {
			watcher.Close()
		}
	}
	this.mu.Unlock()
}

func (this *Registry) Register(ctx context.Context, domain, service, node, addr string, ttl int64) (key string, err error) {
	_, key, err = this.client.Register(ctx, this.BuildPath(domain, service, node), addr, ttl)
	return key, err
}

func (this *Registry) Unregister(ctx context.Context, domain, service, node string) (err error) {
	return this.client.Unregister(ctx, this.BuildPath(domain, service, node))
}

func (this *Registry) BuildPath(domain, service, node string) string {
	return this.buildPath(domain, service, node)
}

func (this *Registry) buildPath(paths ...string) string {
	var nPath = filepath.Join(paths...)

	if len(nPath) > 0 && nPath[0] == '/' {
		nPath = nPath[1:]
	}

	var buf = bytes.NewBufferString(this.scheme)
	buf.WriteString("://")
	buf.WriteString(nPath)

	if len(nPath) > 0 && nPath[len(nPath)-1] != '/' {
		buf.WriteString("/")
	}
	return buf.String()
}
