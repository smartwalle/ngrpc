package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/naming/etcd/resolver"
	"google.golang.org/grpc"
	"time"
)

func main() {
	var builder = resolver.NewBuilder(examples.GetETCDClient())

	var conn = ngrpc.Dial(builder.BuildPath("grpc3", "s2", ""),
		grpc.WithResolvers(builder),
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

	conn.Close()

	examples.Wait()
}
