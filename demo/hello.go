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
	return rsp, nil
}

func (this *HelloService) Stream(s proto.HelloWorld_StreamServer) error {
	for {
		m, err := s.Recv()
		if err != nil {
			return err
		}

		var w = &proto.World{}
		w.Message = fmt.Sprintf("收到来自 %s 的消息", m.Name)
		if err = s.Send(w); err != nil {
			return err
		}
	}
	return nil
}
