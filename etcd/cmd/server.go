package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/grpc4go/etcd/cmd/hw"
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
	fmt.Println(r.Register("game", "user", "node1", addr, 5))

	server := grpc.NewServer()
	hw.RegisterFirstGRPCServer(server, &service{})
	server.Serve(listener)
}

type service struct {
}

func (this *service) FirstCall(ctx context.Context, req *hw.FirstRequest) (*hw.FirstResponse, error) {
	return &hw.FirstResponse{Message: fmt.Sprintf("Hello %s, from xxx %s", req.Name, addr)}, nil
}
