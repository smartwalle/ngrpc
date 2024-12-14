package main

import (
	"context"
	"github.com/smartwalle/net4go"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/naming/etcd/registry"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	var ip, _ = net4go.GetInternalIP()
	listener, err := net.Listen("tcp", ip+":0")
	if err != nil {
		log.Println(context.Background(), err)
		return
	}

	var domain = "grpc"
	var service = "hello"
	var node = "pool"

	var r = registry.NewRegistry(examples.GetETCDClient())
	r.Register(context.Background(), domain, service, node, listener.Addr().String(), 10)

	var server = grpc.NewServer()
	proto.RegisterHelloWorldServer(server, &examples.HelloService{})

	go func() {
		log.Println(context.Background(), "服务地址:", listener.Addr().String())
		err = server.Serve(listener)
		if err != nil {
			log.Println(context.Background(), "启动服务发生错误:", err)
		}
	}()

	examples.Wait()

	// 关闭服务
	server.Stop()

	// 取消注册服务
	r.Unregister(context.Background(), domain, service, node)
}
