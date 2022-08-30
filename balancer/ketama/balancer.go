package ketama

import (
	"github.com/smartwalle/nhash/ketama"
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
	h   func() hash.Hash32
	key string
}

func (this *PickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var picker = &Picker{}
	picker.key = this.key
	picker.selector = ketama.New[balancer.SubConn](ketama.WithSpots(8), ketama.WithHash(this.h))
	for conn, connInfo := range info.ReadySCs {
		picker.selector.Add(connInfo.Address.Addr, conn, 1)
	}
	picker.selector.Prepare()
	return picker
}

type Picker struct {
	selector *ketama.Hash[balancer.SubConn]
	key      string
}

func (this *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var r balancer.PickResult
	var value, _ = info.Ctx.Value(this.key).(string)
	r.SubConn = this.selector.Get(value)
	return r, nil
}
