package ketama

import (
	"github.com/smartwalle/hash4go/ketama"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"hash"
)

const Name = "ngrpc_balancer_ketama"

// New 创建一致性 Hash 负载均衡器
// 使用： grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, ketama.Name())),
func New(key string, h func() hash.Hash32) balancer.Builder {
	var b = base.NewBalancerBuilder(Name, &kPickerBuilder{key: key, h: h}, base.Config{HealthCheck: true})
	balancer.Register(b)
	return b
}

type kPickerBuilder struct {
	key string
	h   func() hash.Hash32
}

func (this *kPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var picker = &kPicker{}
	picker.key = this.key
	picker.selector = ketama.New(8, this.h)
	for conn, cInfo := range info.ReadySCs {
		picker.selector.Add(cInfo.Address.Addr, conn, 1)
	}
	picker.selector.Prepare()
	return picker
}

type kPicker struct {
	key      string
	selector *ketama.Hash
}

func (this *kPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	var r balancer.PickResult
	var value, _ = info.Ctx.Value(this.key).(string)
	r.SubConn, _ = this.selector.Get(value).(balancer.SubConn)
	return r, nil
}
