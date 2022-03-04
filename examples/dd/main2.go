package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/examples"
	"github.com/smartwalle/grpc4go/examples/proto"
	"github.com/smartwalle/grpc4go/registry/etcd"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/xid"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var r = etcd.NewRegistry(examples.GetETCDClient())
	var s, err = grpc4go.NewServer("grpc3", "s2", xid.NewMID().Hex(),
		r,
		grpc4go.WithRegisterTTL(5),
	)
	if err != nil {
		return
	}

	var conn = grpc4go.Dial(r.BuildPath("grpc3", "s1", ""),
		grpc4go.WithPoolSize(6),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	conn.Prepare()

	time.Sleep(time.Second * 1)

	var client = proto.NewHelloWorldClient(conn)

	proto.RegisterHelloWorldServer(s, &Service2{client: client})

	go func() {
		log4go.Println(nil, "服务地址:", s.Addr(), s.Name())
		var err = s.Run()
		if err != nil {
			log4go.Println(nil, "启动服务发生错误:", err)
		}
	}()

	examples.Wait()

	// 关闭服务
	s.Stop()
	r.Close()
}

type Service2 struct {
	proto.UnimplementedHelloWorldServer
	client proto.HelloWorldClient
}

func (this *Service2) Call(ctx context.Context, in *proto.Hello) (*proto.World, error) {
	log4go.Println(ctx, "收到请求", in.Name)

	var rsp = &proto.World{}
	rsp.Message = fmt.Sprintf("收到来自 %s 的消息", in.Name)

	_, err := this.client.Call(ctx, in)
	fmt.Println("req", err)

	return rsp, nil
}

func (this *Service2) Stream(s proto.HelloWorld_StreamServer) error {
	return nil

}
