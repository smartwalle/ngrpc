package main

import (
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/demo"
	"github.com/smartwalle/grpc4go/demo/proto"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/log4go"
)

func main() {
	var r = etcd.NewRegistry(demo.GetETCDClient())
	var server = grpc4go.NewServer("grpc2", "hello", "cmd2", demo.GetIPAddress(), r)

	proto.RegisterHelloWorldServer(server.Server(), &demo.HelloService{})

	go func() {
		log4go.Println("服务地址:", server.Addr(), server.Name())
		var err = server.Run()
		if err != nil {
			log4go.Println("启动服务发生错误:", err)
		}
	}()

	demo.Wait()

	// 关闭服务
	server.Stop()
	r.Close()
}
