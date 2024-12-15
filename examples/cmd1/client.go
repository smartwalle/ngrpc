package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/naming/etcd/resolver"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	var builder = resolver.NewBuilder(examples.GetETCDClient())
	builder.Register()

	var conn = ngrpc.Dial(
		builder.BuildPath("grpc1", "hello", ""),
		ngrpc.WithInsecure(),
		ngrpc.WithTimeout(time.Second*3),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "round_robin"}`)),
	)

	go func() {
		var client = proto.NewHelloWorldClient(conn)

		for i := 0; i < 100; i++ {
			rsp, err := client.Call(context.Background(), &proto.Hello{Name: "Coffee"})
			if err != nil {
				log.Println(context.Background(), err)
				time.Sleep(time.Second * 1)
				continue
			}
			log.Println(context.Background(), i, rsp.Message)

			time.Sleep(time.Second * 1)
		}
	}()

	examples.Wait()

	conn.Close()

	examples.Wait()
}
