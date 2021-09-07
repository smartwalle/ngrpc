package ketama

import (
	"github.com/smartwalle/hash4go/ketama"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const Name = "grpc4go_balancer_ketama"

// New 创建一致性 Hash 负载均衡器
// 使用： grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, ketama.Name())),
func New(key string) balancer.Builder {
	var b = base.NewBalancerBuilder(Name, &kPickerBuilder{key: key}, base.Config{HealthCheck: true})
	balancer.Register(b)
	return b
}

type kPickerBuilder struct {
	key string
}

func (this *kPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var picker = &kPicker{}
	picker.key = this.key
	picker.selector = ketama.New(8, nil)
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
