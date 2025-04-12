package main

import (
	"context"
	"encoding/json"
	"github.com/smartwalle/ngrpc"
	"github.com/smartwalle/ngrpc/examples"
	"github.com/smartwalle/ngrpc/examples/proto"
	"github.com/smartwalle/ngrpc/interceptor/tracing"
	"github.com/smartwalle/ngrpc/naming/etcd/registry"
	"github.com/smartwalle/xid"
	"log"
)

func main() {
	var cfg = examples.Load("./cfg.yaml")
	closer, _ := cfg.InitGlobalTracer("server")

	var r = registry.NewRegistry(examples.GetETCDClient())

	var s, err = ngrpc.NewServer("grpc2", "hello", xid.NewMID().Hex(),
		r,
		ngrpc.WithRegisterTTL(5),
		//wrapper.WithUnaryServer(wrapper.WithWrapper(func(ctx context.Context, md metadata.MD) (context.Context, metadata.MD) {
		//	var logId = md.Get("log-id")[0]
		//	return log.ContextWithId(ctx, logId), md
		//})),
		//wrapper.WithStreamServer(wrapper.WithWrapper(func(ctx context.Context, md metadata.MD) (context.Context, metadata.MD) {
		//	var logId = md.Get("log-id")[0]
		//	return log.ContextWithId(ctx, logId), md
		//})),
		tracing.WithUnaryServer(tracing.WithPayload(true)),
		tracing.WithStreamServer(
			tracing.WithPayload(true),
			tracing.WithStreamPayload(true),
			tracing.WithPayloadMarshal(func(m interface{}) interface{} {
				var data, _ = json.Marshal(m)
				return string(data)
			}),
		),
	)
	if err != nil {
		log.Println(nil, "创建服务发生错误:", err)
		return
	}

	proto.RegisterHelloWorldServer(s, &examples.HelloService{})

	go func() {
		log.Println(nil, "服务地址:", s.Addr(), s.Name())
		var err = s.Start(context.Background())
		if err != nil {
			log.Println(nil, "启动服务发生错误:", err)
		}
	}()

	examples.Wait()

	// 关闭服务
	s.Stop(context.Background())

	closer.Close()
}
