package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/registry/etcd"
	"github.com/smartwalle/xid"
)

func main() {
	var r = etcd.NewRegistry(examples.GetETCDClient())
	var s, err = ngrpc.NewServer("grpc3", "s1", xid.NewMID().Hex(),
		r,
		ngrpc.WithRegisterTTL(5),
	)
	if err != nil {
		return
	}
	proto.RegisterHelloWorldServer(s, &Service1{})

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

type Service1 struct {
	proto.UnimplementedHelloWorldServer
}

func (this *Service1) Call(ctx context.Context, in *proto.Hello) (*proto.World, error) {
	log4go.Println(ctx, "收到请求", in.Name)

	var rsp = &proto.World{}
	rsp.Message = fmt.Sprintf("收到来自 %s 的消息", in.Name)

	return rsp, nil
}

func (this *Service1) Stream(s proto.HelloWorld_StreamServer) error {
	return nil

}
