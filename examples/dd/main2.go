package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/naming/etcd/registry"
	"github.com/smartwalle/ngrpc/naming/etcd/resolver"
	"github.com/smartwalle/xid"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var r = registry.NewRegistry(examples.GetETCDClient())
	var s, err = ngrpc.NewServer("grpc3", "s2", xid.NewMID().Hex(),
		r,
		ngrpc.WithRegisterTTL(5),
	)
	if err != nil {
		return
	}

	var builder = resolver.NewBuilder(examples.GetETCDClient())

	var conn = ngrpc.Dial(builder.BuildPath("grpc3", "s1", ""),
		grpc.WithResolvers(builder),
		ngrpc.WithPoolSize(6),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	conn.Prepare()

	time.Sleep(time.Second * 1)

	var client = proto.NewHelloWorldClient(conn)

	proto.RegisterHelloWorldServer(s, &Service2{client: client})

	go func() {
		log.Println(nil, "服务地址:", s.Addr(), s.Name())
		var err = s.Run()
		if err != nil {
			log.Println(nil, "启动服务发生错误:", err)
		}
	}()

	examples.Wait()

	// 关闭服务
	s.Stop()

	examples.Wait()

}

type Service2 struct {
	proto.UnimplementedHelloWorldServer
	client proto.HelloWorldClient
}

func (this *Service2) Call(ctx context.Context, in *proto.Hello) (*proto.World, error) {
	log.Println(ctx, "收到请求", in.Name)

	var rsp = &proto.World{}
	rsp.Message = fmt.Sprintf("收到来自 %s 的消息", in.Name)

	_, err := this.client.Call(ctx, in)
	fmt.Println("req", err)

	return rsp, nil
}

func (this *Service2) Stream(s proto.HelloWorld_StreamServer) error {
	return nil

}
