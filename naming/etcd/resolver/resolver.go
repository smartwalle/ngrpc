package resolver

import (
	"bytes"
	"context"
	"github.com/smartwalle/ngrpc/naming/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"path/filepath"
)

type Builder struct {
	scheme string
	client *clientv3.Client
}

func NewBuilder(client *clientv3.Client) *Builder {
	return NewBuilderWithScheme(etcd.Scheme, client)
}

func NewBuilderWithScheme(scheme string, client *clientv3.Client) *Builder {
	var b = &Builder{}
	b.scheme = scheme
	b.client = client
	return b
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	var ctx, cancel = context.WithCancel(context.Background())
	var r = &Resolver{
		client: b.client,
		target: target.URL.String(),
		cc:     cc,
		cancel: cancel,
	}
	if err := r.watch(ctx); err != nil {
		return nil, err
	}
	return r, nil
}

func (b *Builder) Scheme() string {
	return b.scheme
}

func (b *Builder) BuildPath(domain, service, node string) string {
	return b.buildPath(domain, service, node)
}

func (b *Builder) buildPath(paths ...string) string {
	var nPath = filepath.Join(paths...)

	if len(nPath) > 0 && nPath[0] == '/' {
		nPath = nPath[1:]
	}

	var buf = bytes.NewBufferString(b.scheme)
	buf.WriteString("://")
	buf.WriteString(nPath)

	if len(nPath) > 0 && nPath[len(nPath)-1] != '/' {
		buf.WriteString("/")
	}
	return buf.String()
}

type Resolver struct {
	client *clientv3.Client
	target string
	cc     resolver.ClientConn
	cancel context.CancelFunc
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (r *Resolver) Close() {
	r.cancel()
}

func (r *Resolver) watch(ctx context.Context) error {
	getRsp, err := r.client.Get(ctx, r.target, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	var eventsChan = make(chan []*clientv3.Event, 1)
	if len(getRsp.Kvs) > 0 {
		var events = make([]*clientv3.Event, 0, len(getRsp.Kvs))
		for _, kv := range getRsp.Kvs {
			events = append(events, &clientv3.Event{Type: clientv3.EventTypePut, Kv: kv})
		}
		eventsChan <- events
	}

	var watchChan = r.client.Watch(ctx, r.target, clientv3.WithRev(getRsp.Header.Revision+1), clientv3.WithPrefix())
	go func() {
		defer close(eventsChan)

		var eventMap = make(map[string]*clientv3.Event)

		for {
			select {
			case <-ctx.Done():
				return
			case watchRsp, ok := <-watchChan:
				if !ok {
					return
				}
				if watchRsp.Err() != nil {
					return
				}
				eventsChan <- watchRsp.Events
			case events, ok := <-eventsChan:
				if !ok {
					return
				}
				for _, event := range events {
					switch event.Type {
					case clientv3.EventTypePut:
						eventMap[string(event.Kv.Key)] = event
					case clientv3.EventTypeDelete:
						delete(eventMap, string(event.Kv.Key))
					default:
					}
				}
				var addrs = make([]resolver.Address, 0, len(eventMap))
				for _, event := range eventMap {
					var addr = resolver.Address{Addr: string(event.Kv.Value)}
					addrs = append(addrs, addr)
				}
				r.cc.UpdateState(resolver.State{Addresses: addrs})
			}
		}
	}()
	return nil
}
