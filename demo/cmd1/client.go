package main

import (
	"context"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/demo"
	"github.com/smartwalle/grpc4go/demo/proto"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/log4go"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var r = etcd.NewRegistry(demo.GetETCDClient())

	var conn = grpc4go.Dial(r.BuildPath("grpc1", "hello", "cmd1"), 10, time.Second*3, grpc.WithInsecure())

	var client = proto.NewHelloWorldClient(conn)

	for i := 0; i < 100; i++ {
		rsp, err := client.Call(context.Background(), &proto.Hello{Name: "Coffee"})
		if err != nil {
			log4go.Println(err)
			return
		}
		log4go.Println(i, rsp.Message)

		time.Sleep(time.Second * 1)
	}

	r.Close()
}
