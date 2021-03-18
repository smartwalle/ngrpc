package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/grpc4go/etcd/cmd/hw"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var config = clientv3.Config{}
	config.Endpoints = []string{"192.168.1.77:2379"}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	etcd.NewRegistry(etcdClient)

	conn, err := grpc.Dial("etcd://game/user", grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		cc := hw.NewFirstGRPCClient(conn)
		time.Sleep(time.Second * 1)
		rsp, err := cc.FirstCall(context.Background(), &hw.FirstRequest{Name: "Yang"})

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("rand", rsp.Message)
	}
}
