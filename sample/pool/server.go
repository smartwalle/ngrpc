package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/sample/hw"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"net"
)

var addr = ":5006"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 初始化 etcd 连接配置文件
	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}

	// 注册服务
	var r = grpc4go.NewETCDResolverWithConfig(config)
	fmt.Println(r.RegisterService("service", "hello", "node2", addr, 5))

	server := grpc.NewServer()
	hw.RegisterFirstGRPCServer(server, &service{})
	server.Serve(listener)
}

type service struct {
}

func (this *service) FirstCall(ctx context.Context, req *hw.FirstRequest) (*hw.FirstResponse, error) {
	return &hw.FirstResponse{Message: fmt.Sprintf("Hello %s, from xxx %s", req.Name, addr)}, nil
}
