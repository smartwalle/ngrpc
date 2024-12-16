package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/net4go"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/naming/etcd/registry"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var ip, _ = net4go.GetInternalIP()
	listener, err := net.Listen("tcp", ip+":0")
	if err != nil {
		log.Println(context.Background(), err)
		return
	}

	var node = fmt.Sprintf("%d-%d", os.Getpid(), time.Now().Unix())

	var r = registry.NewRegistry(examples.GetETCDClient())
	r.Register(context.Background(), "grpc1", "hello", node, listener.Addr().String(), 10)

	var s = grpc.NewServer()
	proto.RegisterHelloWorldServer(s, &examples.HelloService{})
	go func() {
		log.Println(context.Background(), "服务地址:", listener.Addr().String())
		err = s.Serve(listener)
		if err != nil {
			log.Println(context.Background(), "启动服务发生错误:", err)
		}
	}()

	examples.Wait()

	// 关闭服务
	s.Stop()
	// 取消注册服务
	r.Deregister(context.Background(), "grpc1", "hello", node)

	examples.Wait()
}
