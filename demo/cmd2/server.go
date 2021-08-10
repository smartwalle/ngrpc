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
	var s, err = grpc4go.NewServer("grpc2", "hello", "cmd2", r, grpc4go.WithRegisterTTL(5))
	if err != nil {
		log4go.Println("创建服务发生错误:", err)
		return
	}

	proto.RegisterHelloWorldServer(s, &demo.HelloService{})

	go func() {
		log4go.Println("服务地址:", s.Addr(), s.Name())
		var err = s.Run()
		if err != nil {
			log4go.Println("启动服务发生错误:", err)
		}
	}()

	demo.Wait()

	// 关闭服务
	s.Stop()
	r.Close()
}
