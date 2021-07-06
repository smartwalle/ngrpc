package demo

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go/demo/proto"
)

type HelloService struct {
	proto.UnimplementedHelloWorldServer
}

func (this *HelloService) Call(ctx context.Context, in *proto.Hello) (*proto.World, error) {
	var rsp = &proto.World{}
	rsp.Message = fmt.Sprintf("收到来自 %s 的消息", in.Name)
	//time.Sleep(time.Second * 100)
	return rsp, nil
}
