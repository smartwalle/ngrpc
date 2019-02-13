package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/sample/hw"
	"go.etcd.io/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"os"
	"time"
)

func main() {
	// 初始化 etcd 连接配置文件
	var config = clientv3.Config{}
	config.Endpoints = []string{"localhost:2379"}

	// 注册命名解析及服务发现
	var r = grpc4go.NewETCDResolverWithConfig(config)
	resolver.Register(r)

	var h = grpc4go.NewPoolHub()
	h.AddPool("hello", grpc4go.NewPool(r.GetServicePath("service", "hello", ""), 2, 1, grpc.WithInsecure()))

	go req(h.GetPool("hello"), "")
	go req(h.GetPool("hello"), "node1")
	go req(h.GetPool("hello"), "node2")

	select {}
}

func req(p *grpc4go.Pool, node string) {
	for {
		c := p.GetConn(node)

		cc := hw.NewFirstGRPCClient(c)
		rsp, err := cc.FirstCall(context.Background(), &hw.FirstRequest{Name: "Yang"})

		time.Sleep(time.Microsecond * 100)
		p.Release(c)

		if err != nil {
			fmt.Println(node, err)
			os.Exit(-1)
			continue
		}

		fmt.Printf("%s --- %s --- %p --- %s \n", node, c.Target(), *&c, rsp.Message)
	}
}
