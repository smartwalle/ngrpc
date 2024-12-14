module github.com/smartwalle/ngrpc/examples

go 1.22

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/smartwalle/net4go v0.0.51
	github.com/smartwalle/net4go/grpc v0.0.11
	github.com/smartwalle/ngrpc v0.0.0
	github.com/smartwalle/ngrpc/balancer/ketama v0.0.0
	github.com/smartwalle/ngrpc/naming/etcd v0.0.0
	github.com/smartwalle/xid v1.0.6
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	go.etcd.io/etcd/client/v3 v3.5.17
	google.golang.org/grpc v1.69.0
	google.golang.org/protobuf v1.35.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/smartwalle/nhash v0.0.1 // indirect
	github.com/smartwalle/queue v0.0.3 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.etcd.io/etcd/api/v3 v3.5.17 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.17 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace (
	github.com/smartwalle/ngrpc => ../
	github.com/smartwalle/ngrpc/balancer/ketama => ../balancer/ketama
	github.com/smartwalle/ngrpc/naming/etcd => ../naming/etcd
)
