package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WithUnaryServer 服务端普通方法响应追踪
func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		tracer:         opentracing.GlobalTracer(),
		payloadMarshal: defaultPayloadMarshal,
		opName:         defaultOperationName,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerTracing(defaultOption))
}

func unaryServerTracing(opt *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if opt.disable {
			return handler(ctx, req)
		}
		var opName = opt.opName(ctx, info.FullMethod)

		var nCtx, nSpan, err = serverSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Server] %s", opName))
		if err != nil {
			return nil, err
		}

		if opt.payload {
			var md, _ = metadata.FromIncomingContext(ctx)
			traceHeader(nSpan, md)

			nSpan.LogKV("Req", opt.payloadMarshal(req))
		}

		resp, err := handler(nCtx, req)

		if opt.payload {
			if err == nil && resp != nil {
				nSpan.LogKV("Resp", opt.payloadMarshal(resp))
			}
		}

		finish(nSpan, err)
		return resp, err
	}
}
