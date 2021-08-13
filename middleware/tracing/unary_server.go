package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		tracer: opentracing.GlobalTracer(),
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerTracing(defaultOption))
}

func unaryServerTracing(defaultOption *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if defaultOption.disable {
			return handler(ctx, req)
		}

		var nCtx, nSpan, err = serverSpanFromContext(ctx, defaultOption.tracer, fmt.Sprintf("[GRPC Server] %s", info.FullMethod))
		if err != nil {
			return nil, err
		}
		rsp, err := handler(nCtx, req)
		finish(nSpan, err)
		return rsp, err
	}
}
