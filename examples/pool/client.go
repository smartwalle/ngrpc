package main

import (
	"context"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/registry/etcd"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type Conn interface {
	grpc.ClientConnInterface
	Close() error
}

func main() {
	var registry = etcd.NewRegistry(examples.GetETCDClient())

	var domain = "grpc"
	var service = "hello"
	var node = "pool"

	var target = registry.BuildPath(domain, service, node)

	var conn Conn
	conn, _ = grpc.Dial(
		target,
		ngrpc.WithInsecure(),
		ngrpc.WithTimeout(time.Second*3),
	)
	request("未使用连接池：", conn)
	conn.Close()

	conn = ngrpc.Dial(
		target,
		ngrpc.WithInsecure(),
		ngrpc.WithTimeout(time.Second*3),
		ngrpc.WithPoolSize(1),
	)
	request("连接池，数量 1：", conn)
	conn.Close()

	conn = ngrpc.Dial(
		target,
		ngrpc.WithInsecure(),
		ngrpc.WithTimeout(time.Second*3),
		ngrpc.WithPoolSize(10),
	)
	request("连接池，数量 10：", conn)
	conn.Close()

	conn = ngrpc.Dial(
		target,
		ngrpc.WithInsecure(),
		ngrpc.WithTimeout(time.Second*3),
		ngrpc.WithPoolSize(20),
	)
	request("连接池，数量 20：", conn)
	conn.Close()

	select {}

	registry.Close()
}

func request(comment string, conn grpc.ClientConnInterface) {
	var client = proto.NewHelloWorldClient(conn)
	var wait = sync.WaitGroup{}
	var begin = time.Now()
	wait.Add(1000000)
	for i := 0; i < 1000000; i++ {
		go func() {
			_, err := client.Call(context.Background(), &proto.Hello{Name: "Coffee"})
			wait.Done()
			if err != nil {
				log.Println(context.Background(), err)
				return
			}
		}()
	}
	wait.Wait()
	var diff = time.Now().Sub(begin)
	log.Println(comment, "耗时:", diff)
}
