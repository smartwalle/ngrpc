module github.com/smartwalle/grpc4go/examples

go 1.12

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/smartwalle/etcd4go v1.0.9 // indirect
	github.com/smartwalle/grpc4go v0.0.0
	github.com/smartwalle/grpc4go/balancer/ketama v0.0.0
	github.com/smartwalle/grpc4go/registry/etcd v0.0.0
	github.com/smartwalle/jaeger4go v1.0.0 // indirect
	github.com/smartwalle/log4go v1.0.4
	github.com/smartwalle/net4go v0.0.44
	github.com/smartwalle/net4go/grpc v0.0.0
	github.com/smartwalle/xid v1.0.6 // indirect
	github.com/uber/jaeger-client-go v2.16.0+incompatible // indirect
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634 // indirect
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/smartwalle/grpc4go => ../
	github.com/smartwalle/grpc4go/balancer/ketama => ../balancer/ketama
	github.com/smartwalle/grpc4go/registry/etcd => ../registry/etcd
	github.com/smartwalle/log4go => /Users/yang/Desktop/smartwalle/log4go
	github.com/smartwalle/net4go => /Users/yang/Desktop/smartwalle/net4go
	github.com/smartwalle/net4go/grpc => /Users/yang/Desktop/smartwalle/net4go/grpc
)
