package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go/cmd/hw"
	"github.com/smartwalle/grpc4go/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
)

var addr = ":5005"

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var config = clientv3.Config{}
	config.Endpoints = []string{"192.168.1.77:2379"}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	var r = etcd.NewRegistry(etcdClient)
	r.Register(context.Background(), "game", "user", "node1", addr, 5)

	server := grpc.NewServer()
	hw.RegisterFirstGRPCServer(server, &service{})
	server.Serve(listener)
}

type service struct {
}

func (this *service) FirstCall(ctx context.Context, req *hw.FirstRequest) (*hw.FirstResponse, error) {
	return &hw.FirstResponse{Message: fmt.Sprintf("Hello %s, from xxx %s", req.Name, addr)}, nil
}
