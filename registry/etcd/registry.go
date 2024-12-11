package etcd

import (
	"bytes"
	"context"
	"github.com/smartwalle/netcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
	"sync"
)

const (
	kDefaultScheme = "etcd"
)

type Registry struct {
	client   *netcd.Client
	mu       *sync.Mutex
	watchers map[string]*netcd.Watcher
	leaseIds map[string][]int64
	scheme   string
}

func NewRegistry(client *clientv3.Client) *Registry {
	return NewRegistryWithScheme(kDefaultScheme, client)
}

func NewRegistryWithScheme(scheme string, client *clientv3.Client) *Registry {
	var r = &Registry{scheme: scheme, client: netcd.NewClient(client)}
	r.mu = &sync.Mutex{}
	r.watchers = make(map[string]*netcd.Watcher)
	r.leaseIds = make(map[string][]int64)
	resolver.Register(r)
	return r
}

func (r *Registry) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	//var key = target.Scheme + "://" + filepath.Join(target.Authority, target.Endpoint)
	var key = target.URL.String()
	var watcher = r.client.Watch(context.Background(), key, r.watch(cc), clientv3.WithPrefix())
	r.mu.Lock()
	r.update(cc, watcher.Values())
	r.watchers[key] = watcher
	r.mu.Unlock()
	return r, nil
}

func (r *Registry) watch(cc resolver.ClientConn) func(watcher *netcd.Watcher, event, key, path string, value []byte) {
	return func(watcher *netcd.Watcher, event, key, path string, value []byte) {
		var values = watcher.Values()
		r.update(cc, values)
	}
}

func (r *Registry) update(cc resolver.ClientConn, values map[string][]byte) {
	var addresses = make([]resolver.Address, 0, len(values))
	for _, value := range values {
		var addr = resolver.Address{Addr: string(value)}
		addresses = append(addresses, addr)
	}
	cc.UpdateState(resolver.State{Addresses: addresses})
}

func (r *Registry) Scheme() string {
	return r.scheme
}

func (r *Registry) ResolveNow(options resolver.ResolveNowOptions) {
}

func (r *Registry) Close() {
	r.mu.Lock()
	for _, watcher := range r.watchers {
		if watcher != nil {
			watcher.Close()
		}
	}
	r.mu.Unlock()
}

func (r *Registry) Register(ctx context.Context, domain, service, node, addr string, ttl int64) (key string, err error) {
	var leaseId = int64(0)
	key = r.BuildPath(domain, service, node)
	leaseId, key, err = r.client.Register(ctx, key, addr, ttl)
	if leaseId > 0 {
		r.mu.Lock()
		r.leaseIds[key] = append(r.leaseIds[key], leaseId)
		r.mu.Unlock()
	}
	return key, err
}

func (r *Registry) Unregister(ctx context.Context, domain, service, node string) (err error) {
	var key = r.BuildPath(domain, service, node)
	r.mu.Lock()
	var leaseIds = r.leaseIds[key]
	delete(r.leaseIds, key)
	r.mu.Unlock()
	for _, leaseId := range leaseIds {
		r.client.Revoke(ctx, leaseId)
	}
	return r.client.Unregister(ctx, key)
}

func (r *Registry) BuildPath(domain, service, node string) string {
	return r.buildPath(domain, service, node)
}

func (r *Registry) buildPath(paths ...string) string {
	var nPath = filepath.Join(paths...)

	if len(nPath) > 0 && nPath[0] == '/' {
		nPath = nPath[1:]
	}

	var buf = bytes.NewBufferString(r.scheme)
	buf.WriteString("://")
	buf.WriteString(nPath)

	if len(nPath) > 0 && nPath[len(nPath)-1] != '/' {
		buf.WriteString("/")
	}
	return buf.String()
}
