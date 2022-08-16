package ketama

import (
	"github.com/smartwalle/hash4go/ketama"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"hash"
)

const Name = "ngrpc_balancer_ketama"

// New 创建一致性 Hash 负载均衡器
//
// 使用： grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, ketama.Name))
//
func New(key string, h func() hash.Hash32) balancer.Builder {
	var b = base.NewBalancerBuilder(Name, &PickerBuilder{key: key, h: h}, base.Config{HealthCheck: true})
	balancer.Register(b)
	return b
}

type PickerBuilder struct {
	key string
	h   func() hash.Hash32
}

func (this *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var p = &picker{}
	p.key = this.key
	p.selector = ketama.New(8, this.h)
	for conn, connInfo := range info.ReadySCs {
		p.selector.Add(connInfo.Address.Addr, conn, 1)
	}
	p.selector.Prepare()
	return p
}

type picker struct {
	key      string
	selector *ketama.Hash
}

func (this *picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var r balancer.PickResult
	var value, _ = info.Ctx.Value(this.key).(string)
	r.SubConn, _ = this.selector.Get(value).(balancer.SubConn)
	return r, nil
}
