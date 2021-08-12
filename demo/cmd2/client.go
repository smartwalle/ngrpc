package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/demo"
	"github.com/smartwalle/grpc4go/demo/proto"
	"github.com/smartwalle/grpc4go/etcd"
	"github.com/smartwalle/grpc4go/middleware/logging"
	"github.com/smartwalle/grpc4go/middleware/timeout"
	"github.com/smartwalle/log4go"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var r = etcd.NewRegistry(demo.GetETCDClient())

	log4go.SharedInstance().DisablePath()

	var conn = grpc4go.Dial(r.BuildPath("grpc2", "hello", "cmd2"),
		grpc.WithInsecure(),
		//grpc4go.WithTimeout(time.Second*3),
		timeout.WithUnaryClient(),
		logging.WithUnaryClient(logging.WithLogger(log4go.SharedInstance()), logging.WithPayload(true)),
		logging.WithStreamClient(logging.WithLogger(log4go.SharedInstance()), logging.WithPayload(true)),
	)

	var client = proto.NewHelloWorldClient(conn)

	//go doStream(client)
	go doUnary(client)

	select {}
	r.Close()
}

func doStream(client proto.HelloWorldClient) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			var header = grpc4go.NewHeader()
			header.Set("user-id", "1")
			header.Set("user-id", "2")

			stream, _ := client.Stream(header.Context(context.Background()))

			go func() {
				for {
					_, err := stream.Recv()
					if err != nil {
						return
					}
				}
			}()

			for {
				if err := stream.Send(&proto.Hello{Name: fmt.Sprintf("Stream %d", i)}); err != nil {
					return
				}
				time.Sleep(time.Second * 1)
			}
		}(i)
	}
}

func doUnary(client proto.HelloWorldClient) {
	for {

		var header = grpc4go.NewHeader()
		header.Set("user-id", "1")
		header.Set("user-id", "2")

		_, _ = client.Call(header.Context(context.Background()), &proto.Hello{Name: "Coffee"})

		//if err != nil {
		//	log4go.Println(err, rsp)
		//}
		time.Sleep(time.Second * 1)
	}
}
