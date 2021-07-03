module github.com/smartwalle/grpc4go/demo

go 1.12

require (
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/smartwalle/grpc4go v0.0.0
	github.com/smartwalle/grpc4go/etcd v0.0.0
	github.com/smartwalle/log4go v1.0.4
	github.com/smartwalle/net4go v0.0.39
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634 // indirect
	google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884 // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/smartwalle/grpc4go => ../
	github.com/smartwalle/grpc4go/etcd => ../etcd
)
