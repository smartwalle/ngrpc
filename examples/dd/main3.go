package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/registry/etcd"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var r = etcd.NewRegistry(examples.GetETCDClient())

	var conn = ngrpc.Dial(r.BuildPath("grpc3", "s2", ""),
		ngrpc.WithPoolSize(3),
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
