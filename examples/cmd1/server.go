package main

import (
	"context"
	"github.com/smartwalle/grpc4go/examples"
	"github.com/smartwalle/grpc4go/examples/proto"
	"github.com/smartwalle/grpc4go/registry/etcd"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/net4go"
	"google.golang.org/grpc"
	"net"
)

func main() {
	var ip, _ = net4go.GetInternalIP()
	listener, err := net.Listen("tcp", ip+":0")
	if err != nil {
		log4go.Println(context.Background(), err)
		return
	}

	var r = etcd.NewRegistry(examples.GetETCDClient())
	r.Register(context.Background(), "grpc1", "hello", "cmd1", listener.Addr().String(), 10)

	var s = grpc.NewServer()
	proto.RegisterHelloWorldServer(s, &examples.HelloService{})
	go func() {
		log4go.Println(context.Background(), "服务地址:", listener.Addr().String())
		err = s.Serve(listener)
		if err != nil {
			log4go.Println(context.Background(), "启动服务发生错误:", err)
		}
	}()

	examples.Wait()

	// 关闭服务
	s.Stop()
	// 取消注册服务
	r.Unregister(context.Background(), "grpc1", "hello", "cmd1")
	r.Close()
}
