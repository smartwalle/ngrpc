package main

import (
	"context"
	"github.com/smartwalle/grpc4go/demo"
	"github.com/smartwalle/grpc4go/demo/proto"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/net4go"
	"google.golang.org/grpc"
	"net"
)

func main() {
	var ip, _ = net4go.GetInternalIP()
	listener, err := net.Listen("tcp", ip+":0")
	if err != nil {
		log4go.Println(err)
		return
	}

	var r = etcd.NewRegistry(demo.GetETCDClient())
	r.Register(context.Background(), "grpc1", "hello", "cmd1", listener.Addr().String(), 10)

	var server = grpc.NewServer()
	proto.RegisterHelloWorldServer(server, &demo.HelloService{})
	go func() {
		log4go.Println("服务地址:", listener.Addr().String())
		err = server.Serve(listener)
		if err != nil {
			log4go.Println("启动服务发生错误:", err)
		}
	}()

	demo.Wait()

	// 关闭服务
	server.Stop()
	// 取消注册服务
	r.Deregister(context.Background(), "grpc1", "hello", "cmd1")
	r.Close()
}
