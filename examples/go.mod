module github.com/smartwalle/grpc4go/examples

go 1.12

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/smartwalle/grpc4go v0.0.0
	github.com/smartwalle/grpc4go/balancer/ketama v0.0.0
	github.com/smartwalle/grpc4go/registry/etcd v0.0.0
	github.com/smartwalle/log4go v1.0.4
	github.com/smartwalle/net4go v0.0.50
	github.com/smartwalle/net4go/grpc v0.0.10
	github.com/smartwalle/xid v1.0.6
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.etcd.io/etcd/client/v3 v3.5.4
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/smartwalle/grpc4go => ../
	github.com/smartwalle/grpc4go/balancer/ketama => ../balancer/ketama
	github.com/smartwalle/grpc4go/registry/etcd => ../registry/etcd
)
