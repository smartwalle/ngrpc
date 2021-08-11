package main

import (
	"context"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/demo"
	"github.com/smartwalle/grpc4go/demo/proto"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/grpc4go/middleware/logging"
	"github.com/smartwalle/log4go"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var r = etcd.NewRegistry(demo.GetETCDClient())

	var conn = grpc4go.Dial(r.BuildPath("grpc2", "hello", "cmd2"), grpc.WithInsecure(), grpc4go.WithTimeout(time.Second*3), logging.WithUnaryCall(logging.WithLogger(log4go.SharedInstance())))

	var client = proto.NewHelloWorldClient(conn)

	for i := 0; i < 100; i++ {
		_, err := client.Call(context.Background(), &proto.Hello{Name: "Coffee"})
		if err != nil {
			log4go.Println(err)
		}
		time.Sleep(time.Second * 1)
	}

	r.Close()
}
