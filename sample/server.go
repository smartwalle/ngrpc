package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/etcd4go"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/sample/hw"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"net"
)

var addr = ":5004"

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
	var c, _ = etcd4go.NewClient(config)
	var r = grpc4go.NewETCDResolver(c)
	fmt.Println(r.RegisterService("mine/my_service", "node1", addr, 5))

	server := grpc.NewServer()
	hw.RegisterFirstGRPCServer(server, &service{})
	server.Serve(listener)
}

type service struct {
}

func (this *service) FirstCall(ctx context.Context, req *hw.FirstRequest) (*hw.FirstResponse, error) {
	return &hw.FirstResponse{Message: fmt.Sprintf("Hello %s, from %s", req.Name, addr)}, nil
}
