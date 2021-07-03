package main

import (
	"context"
	"github.com/smartwalle/grpc4go/demo"
	"github.com/smartwalle/grpc4go/demo/proto"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/net4go"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var ip, _ = net4go.GetInternalIP()
	listener, err := net.Listen("tcp", ip+":0")
	if err != nil {
		log4go.Println(err)
		return
	}

	var config = clientv3.Config{}
	config.Endpoints = demo.EDCDEndPoints
	etcdClient, err := clientv3.New(config)
	if err != nil {
		log4go.Println(err)
		return
	}

	var r = etcd.NewRegistry(etcdClient)
	r.Register(context.Background(), "grpc1", "hello", "cmd1", listener.Addr().String(), 10)

	log4go.Println("服务地址:", listener.Addr().String())

	server := grpc.NewServer()
	proto.RegisterHelloWorldServer(server, &demo.HelloService{})
	go func() {
		err = server.Serve(listener)
		if err != nil {
			log4go.Println("启动服务发生错误:", err)
		}
	}()

	var c = make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
MainLoop:
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			break MainLoop
		}
	}
	// 取消注册服务
	r.Deregister(context.Background(), "grpc", "hello", "normal")
	r.Close()
}
