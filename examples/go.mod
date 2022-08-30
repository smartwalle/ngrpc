module github.com/smartwalle/ngrpc/examples

go 1.12

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/smartwalle/ngrpc v0.0.0
	github.com/smartwalle/ngrpc/balancer/ketama v0.0.0
	github.com/smartwalle/ngrpc/registry/etcd v0.0.0
	github.com/smartwalle/log4go v1.0.4
	github.com/smartwalle/net4go v0.0.51
	github.com/smartwalle/net4go/grpc v0.0.11
	github.com/smartwalle/xid v1.0.6
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.etcd.io/etcd/client/v3 v3.5.4
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/smartwalle/ngrpc => ../
	github.com/smartwalle/ngrpc/balancer/ketama => ../balancer/ketama
	github.com/smartwalle/ngrpc/registry/etcd => ../registry/etcd
)
