package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/examples"
	"github.com/smartwalle/grpc4go/examples/proto"
	"github.com/smartwalle/grpc4go/registry/etcd"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var r = etcd.NewRegistry(examples.GetETCDClient())

	var conn = grpc4go.Dial(r.BuildPath("grpc3", "s2", ""),
		grpc4go.WithPoolSize(3),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	conn.Prepare()

	time.Sleep(time.Second * 1)

	var client = proto.NewHelloWorldClient(conn)

	var in = &proto.Hello{}
	in.Name = "xxx"

	_, err := client.Call(context.Background(), in)
	fmt.Println("req", err)

	examples.Wait()

	r.Close()
}
