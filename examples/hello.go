package examples

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/smartwalle/net4go"
	ngrpc "github.com/smartwalle/net4go/grpc"
	"github.com/smartwalle/ngrpc/examples/proto"
	"sync"
)

type HelloService struct {
	proto.UnimplementedHelloWorldServer
}

func (hs *HelloService) Call(ctx context.Context, in *proto.Hello) (*proto.World, error) {
	//log.Println("收到请求", in.Name)

	var rsp = &proto.World{}
	rsp.Message = fmt.Sprintf("收到来自 %s 的消息", in.Name)

	span, _ := opentracing.StartSpanFromContext(ctx, "sss")
	span.LogKV("s1", "s1")
	span.Finish()

	return rsp, nil
}

func (hs *HelloService) Stream(s proto.HelloWorld_StreamServer) error {
	fmt.Println("begin...")

	var ns = NewStream(s)

	var nSess = ngrpc.NewSession(ns, hs)

	fmt.Println(nSess.Conn())

	var err = ns.Wait()

	fmt.Println("end...", err)

	//for {
	//	m, err := s.Recv()
	//	if err != nil {
	//		return err
	//	}
	//
	//	log.Println(s.Context(), "收到流消息")
	//
	//	var w = &proto.World{}
	//	w.Message = fmt.Sprintf("收到来自 %s 的消息", m.Name)
	//	if err = s.Send(w); err != nil {
	//		return err
	//	}
	//}
	return err
}

func (hs *HelloService) OnMessage(sess net4go.Session, p net4go.Packet) {
	fmt.Println("OnMessage", p)

	var h = p.(*proto.Hello)

	var w = &proto.World{}
	w.Message = fmt.Sprintf("收到来自 %s 的消息", h.Name)
	if err := sess.WritePacket(w); err != nil {
	}
}

func (hs *HelloService) OnClose(sess net4go.Session, err error) {
	fmt.Println("OnClose", err)
}

type Stream struct {
	s        proto.HelloWorld_StreamServer
	done     chan struct{}
	doneOnce *sync.Once
	err      error
}

func NewStream(s proto.HelloWorld_StreamServer) *Stream {
	var ns = &Stream{}
	ns.s = s
	ns.done = make(chan struct{})
	ns.doneOnce = &sync.Once{}
	return ns
}

func (s *Stream) SendPacket(packet net4go.Packet) error {
	return s.s.Send(packet.(*proto.World))
}

func (s *Stream) RecvPacket() (net4go.Packet, error) {
	return s.s.Recv()
}

func (s *Stream) OnClose(err error) {
	fmt.Println("stream close")
	s.doneOnce.Do(func() {
		close(s.done)
		s.err = err
	})
}

func (s *Stream) Wait() error {
	<-s.done
	return s.err
}
