package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/etcd4go"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/sample/hw"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"time"
)

func main() {
	// 初始化 etcd 连接配置文件
	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}

	// 注册命名解析及服务发现
	var c, _ = etcd4go.NewClient(config)
	var r = grpc4go.NewETCDResolver(c)
	resolver.Register(r)

	var h = grpc4go.NewPoolHub()
	h.AddPool("hello", grpc4go.NewPool(r.GetServicePath("service", "hello", ""), 2, 1, grpc.WithInsecure()))

	for {
		var p = h.GetPool("hello")

		c := p.GetConn()

		fmt.Printf("%p \n", *&c)

		cc := hw.NewFirstGRPCClient(c)
		rsp, err := cc.FirstCall(context.Background(), &hw.FirstRequest{Name: "Yang"})

		time.Sleep(time.Second * 1)
		p.Release(c)

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("rand", rsp.Message)
	}
}
