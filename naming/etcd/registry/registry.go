package registry

import (
	"context"
	"github.com/smartwalle/ngrpc/naming/etcd/internal"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type Registry struct {
	scheme string
	client *clientv3.Client
	mu     *sync.Mutex
	leases map[string]*internal.Lease
}

func NewRegistry(client *clientv3.Client) *Registry {
	return NewRegistryWithScheme(internal.Scheme, client)
}

func NewRegistryWithScheme(scheme string, client *clientv3.Client) *Registry {
	var r = &Registry{}
	r.scheme = scheme
	r.client = client
	r.mu = &sync.Mutex{}
	r.leases = make(map[string]*internal.Lease)
	return r
}

func (r *Registry) Scheme() string {
	return r.scheme
}

func (r *Registry) Register(ctx context.Context, domain, service, node, addr string, ttl int64) (string, error) {
	var key = r.BuildPath(domain, service, node)

	var lease, err = internal.NewLease(r.client, ttl)
	if err != nil {
		return "", err
	}
	r.mu.Lock()
	r.leases[key] = lease
	r.mu.Unlock()

	if _, err = r.client.Put(ctx, key, addr, clientv3.WithLease(lease.ID())); err != nil {
		return "", err
	}
	return key, nil
}

func (r *Registry) Unregister(ctx context.Context, domain, service, node string) error {
	var key = r.BuildPath(domain, service, node)

	r.mu.Lock()
	var lease = r.leases[key]
	delete(r.leases, key)
	r.mu.Unlock()

	if _, err := r.client.Delete(ctx, key); err != nil {
		return err
	}
	return lease.Revoke(ctx)
}

func (r *Registry) BuildPath(domain, service, node string) string {
	return internal.BuildPath(r.scheme, domain, service, node)
}
