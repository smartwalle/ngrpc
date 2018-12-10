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

	for {
		conn, err := grpc.Dial("etcd://service/hello", grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			return
		}

		cc := hw.NewFirstGRPCClient(conn)
		time.Sleep(time.Second * 1)
		rsp, err := cc.FirstCall(context.Background(), &hw.FirstRequest{Name: "Yang"})
		conn.Close()

		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("rand", rsp.Message)
	}
}
