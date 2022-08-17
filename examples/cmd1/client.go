package main

import (
	"context"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/registry/etcd"
	"time"
)

func main() {
	var r = etcd.NewRegistry(examples.GetETCDClient())

	var conn = ngrpc.Dial(r.BuildPath("grpc1", "hello", "cmd1"), ngrpc.WithInsecure(), ngrpc.WithTimeout(time.Second*3))

	var client = proto.NewHelloWorldClient(conn)

	for i := 0; i < 100; i++ {
		rsp, err := client.Call(context.Background(), &proto.Hello{Name: "Coffee"})
		if err != nil {
			log4go.Println(context.Background(), err)
			return
		}
		log4go.Println(context.Background(), i, rsp.Message)

		time.Sleep(time.Second * 1)
	}

	r.Close()
}
