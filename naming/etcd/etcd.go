package etcd

import (
	"github.com/smartwalle/ngrpc/naming/etcd/internal"
	"github.com/smartwalle/ngrpc/naming/etcd/registry"
	"github.com/smartwalle/ngrpc/naming/etcd/resolver"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Registry struct {
	*registry.Registry
	builder *resolver.Builder
}

func New(client *clientv3.Client) *Registry {
	var r = &Registry{}
	r.Registry = registry.NewRegistryWithScheme(internal.Scheme, client)
	r.builder = resolver.NewBuilderWithScheme(internal.Scheme, client)
	return r
}

func NewWithScheme(scheme string, client *clientv3.Client) *Registry {
	var r = &Registry{}
	r.Registry = registry.NewRegistryWithScheme(scheme, client)
	r.builder = resolver.NewBuilderWithScheme(scheme, client)
	return r
}

func (r *Registry) ResolverBuilder() *resolver.Builder {
	return r.builder
}

func (r *Registry) RegisterResolver() {
	r.builder.Register()
}
