package main

import (
	"context"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/examples"
	"github.com/smartwalle/grpc4go/examples/proto"
	"github.com/smartwalle/grpc4go/registry/etcd"
	"github.com/smartwalle/log4go"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var r = etcd.NewRegistry(examples.GetETCDClient())

	var conn = grpc4go.Dial(r.BuildPath("grpc1", "hello", "cmd1"), grpc.WithInsecure(), grpc4go.WithTimeout(time.Second*3))

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
