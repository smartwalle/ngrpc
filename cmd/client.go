package main

import (
	"context"
	"fmt"
	"github.com/smartwalle/grpc4go"
	"github.com/smartwalle/grpc4go/cmd/hw"
	"github.com/smartwalle/grpc4go/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"sync"
	"time"
)

func main() {
	var config = clientv3.Config{}
	config.Endpoints = []string{"192.168.1.77:2379"}
	etcdClient, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	etcd.NewRegistry(etcdClient)

	var conn = grpc4go.Dial("etcd://game/user", 10, time.Second*3, grpc.WithInsecure())

	fmt.Println("ready...")

	var now = time.Now()
	var wg = &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for x := 0; x < 10000; x++ {
				cc := hw.NewFirstGRPCClient(conn)
				if _, err = cc.FirstCall(context.Background(), &hw.FirstRequest{Name: "Yang"}); err != nil {
					fmt.Println("call error:", err)
					continue
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("耗时：", time.Now().Sub(now))
	select {}
}
