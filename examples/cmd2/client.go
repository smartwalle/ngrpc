package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/smartwalle/log4go"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/balancer/ketama"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/interceptor/tracing"
	"github.com/smartwalle/ngrpc/interceptor/wrapper"
	"github.com/smartwalle/ngrpc/registry/etcd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

func main() {
	var cfg = examples.Load("./cfg.yaml")
	closer, _ := cfg.InitGlobalTracer("client")

	var bb = ketama.New("player_id", nil)

	var r = etcd.NewRegistry(examples.GetETCDClient())

	var conn = ngrpc.Dial(r.BuildPath("grpc2", "hello", ""),
		ngrpc.WithPoolSize(3),
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, bb.Name())),
		grpc.WithInsecure(),
		//ngrpc.WithTimeout(time.Second*3),
		//timeout.WithUnaryClient(),

		wrapper.WithUnaryClient(wrapper.WithWrapper(func(ctx context.Context, md metadata.MD) (context.Context, metadata.MD) {
			return ctx, md
		})),
		wrapper.WithStreamClient(wrapper.WithWrapper(func(ctx context.Context, md metadata.MD) (context.Context, metadata.MD) {
			return ctx, md
		})),
		tracing.WithUnaryClient(tracing.WithPayload(true), tracing.WithPayloadMarshal(func(m interface{}) interface{} {
			var data, _ = json.Marshal(m)
			return string(data)
		})),
		tracing.WithStreamClient(tracing.WithPayload(true), tracing.WithStreamPayload(true), tracing.WithPayloadMarshal(func(m interface{}) interface{} {
			var data, _ = json.Marshal(m)
			return string(data)
		})),
	)
	conn.Prepare()
	_ = conn

	time.Sleep(time.Second * 1)

	var client = proto.NewHelloWorldClient(conn)

	fmt.Println("wait...")
	//time.Sleep(time.Second * 2)
	fmt.Println("begin...")

	//go doStream(client)
	go doUnary(client, "8")
	//go doUnary(client, "9")
	//go doUnary(client, "19")
	//go doUnary(client, "20")
	//go doUnary(client, "100")

	//go func() {
	//	time.Sleep(time.Second * 5)
	//
	//	fmt.Println("======")
	//
	//	var conn2, _ = grpc.Dial(
	//		r.BuildPath("grpc2", "hello", ""),
	//		grpc.WithInsecure(),
	//		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, bb.Name())),
	//		wrapper.WithUnaryClient(wrapper.WithWrapper(func(ctx context.Context, md metadata.MD) (context.Context, metadata.MD) {
	//			var logId = log4go.MustGetId(ctx)
	//			md.Set("log-id", logId)
	//			return ctx, md
	//		})),
	//	)
	//	_ = conn2
	//}()

	examples.Wait()

	r.Close()
	closer.Close()
}

func doStream(client proto.HelloWorldClient) {
	for i := 0; i < 10; i++ {
		go func(i int) {
			var header = ngrpc.NewHeader()
			header.Set("user-id", "1")
			header.Set("user-id", "2")

			var ctx = context.Background()

			stream, _ := client.Stream(header.Context(ctx), tracing.WithOperationName(func(ctx context.Context, method string) string {
				return "wtf"
			}))

			go func() {
				for {
					_, err := stream.Recv()
					if err != nil {
						return
					}
				}
			}()

			var c = 0
			for {
				log4go.Println(stream.Context(), "发送流消息")

				if err := stream.Send(&proto.Hello{Name: fmt.Sprintf("Stream %d-%d", i, c)}); err != nil {
					return
				}

				c++

				time.Sleep(time.Second * 1)
			}
		}(i)
	}
}

func doUnary(client proto.HelloWorldClient, id string) {
	var i = 0
	for {
		var ctx = context.Background()

		span, ctx := opentracing.StartSpanFromContext(ctx, "s1-call")
		span.LogKV("s1-call-key", "s1-call-value")
		span.Finish()

		var header = ngrpc.NewHeader()
		header.Set("user-id", "1")
		header.Set("user-id", "2")

		log4go.Println(ctx, "开始请求", id)
		ctx = context.WithValue(ctx, "player_id", id)
		client.Call(header.Context(ctx), &proto.Hello{Name: fmt.Sprintf("Coffee " + id)})
		//log4go.Println(ctx, "请求完成", err)

		span, _ = opentracing.StartSpanFromContext(ctx, "xxxxx")
		span.LogKV("ddd", "ww")
		span.Finish()

		//if err != nil {
		//	log4go.Println(err, rsp)
		//}

		time.Sleep(time.Second * 1)
		i++
	}
}
