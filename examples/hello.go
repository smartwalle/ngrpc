package examples

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/smartwalle/grpc4go/examples/proto"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/net4go"
	ngrpc "github.com/smartwalle/net4go/grpc"
	"sync"
)

type HelloService struct {
	proto.UnimplementedHelloWorldServer
}

func (this *HelloService) Call(ctx context.Context, in *proto.Hello) (*proto.World, error) {
	log4go.Println(ctx, "收到请求", in.Name)

	var rsp = &proto.World{}
	rsp.Message = fmt.Sprintf("收到来自 %s 的消息", in.Name)

	span, _ := opentracing.StartSpanFromContext(ctx, "sss")
	span.LogKV("s1", "s1")
	span.Finish()

	return rsp, nil
}

func (this *HelloService) Stream(s proto.HelloWorld_StreamServer) error {
	fmt.Println("begin...")

	var ns = NewStream(s)

	var nSess = ngrpc.NewSession(ns, this)

	fmt.Println(nSess.Conn())

	var err = ns.Wait()

	fmt.Println("end...", err)

	//for {
	//	m, err := s.Recv()
	//	if err != nil {
	//		return err
	//	}
	//
	//	log4go.Println(s.Context(), "收到流消息")
	//
	//	var w = &proto.World{}
	//	w.Message = fmt.Sprintf("收到来自 %s 的消息", m.Name)
	//	if err = s.Send(w); err != nil {
	//		return err
	//	}
	//}
	return err
}

func (this *HelloService) OnMessage(sess net4go.Session, p net4go.Packet) {
	fmt.Println("OnMessage", p)

	var h = p.(*proto.Hello)

	var w = &proto.World{}
	w.Message = fmt.Sprintf("收到来自 %s 的消息", h.Name)
	if err := sess.WritePacket(w); err != nil {
	}
}

func (this *HelloService) OnClose(sess net4go.Session, err error) {
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

func (this *Stream) SendPacket(packet net4go.Packet) error {
	return this.s.Send(packet.(*proto.World))
}

func (this *Stream) RecvPacket() (net4go.Packet, error) {
	return this.s.Recv()
}

func (this *Stream) OnClose(err error) {
	fmt.Println("stream close")
	this.doneOnce.Do(func() {
		close(this.done)
		this.err = err
	})
}

func (this *Stream) Wait() error {
	<-this.done
	return this.err
}
