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
	var defaultOptions = &options{
		tracer:         opentracing.GlobalTracer(),
		payloadMarshal: defaultPayloadMarshal,
		opName:         defaultOperationName,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.ChainUnaryInterceptor(unaryServerTracing(defaultOptions))
}

func unaryServerTracing(opts *options) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if opts.disable {
			return handler(ctx, req)
		}
		var opName = opts.opName(ctx, info.FullMethod)

		var nCtx, nSpan, err = serverSpanFromContext(ctx, opts.tracer, fmt.Sprintf("[GRPC Server] %s", opName))
		if err != nil {
			return nil, err
		}

		if opts.payload {
			var md, _ = metadata.FromIncomingContext(ctx)
			traceHeader(nSpan, md)

			nSpan.LogKV("Req", opts.payloadMarshal(req))
		}

		resp, err := handler(nCtx, req)

		if opts.payload {
			if err == nil && resp != nil {
				nSpan.LogKV("Resp", opts.payloadMarshal(resp))
			}
		}

		finish(nSpan, err)
		return resp, err
	}
}
