module github.com/smartwalle/grpc4go/cmd

go 1.12

require (
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/smartwalle/grpc4go v0.0.0
	github.com/smartwalle/grpc4go/etcd v0.0.0
	go.etcd.io/etcd/client/v3 v3.5.0-alpha.0 // indirect
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	google.golang.org/grpc v1.32.0 // indirect
)

replace (
	github.com/smartwalle/grpc4go => ../
	github.com/smartwalle/grpc4go/etcd => ../etcd
)