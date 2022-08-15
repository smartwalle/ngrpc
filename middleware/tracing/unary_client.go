package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WithUnaryClient 客户端普通方法调用追踪
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		tracer:         opentracing.GlobalTracer(),
		payloadMarshal: defaultPayloadMarshal,
		opName:         defaultOperationName,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientTracing(defaultOptions))
}

func unaryClientTracing(defaultOptions *options) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)
		if opt.disable {
			return invoker(ctx, method, req, reply, cc, grpcOpts...)
		}
		var opName = opt.opName(ctx, method)

		var nCtx, nSpan, err = clientSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Client] %s", opName))
		if err != nil {
			return err
		}

		if opt.payload {
			var md, _ = metadata.FromOutgoingContext(ctx)
			traceHeader(nSpan, md)

			nSpan.LogKV("Req", opt.payloadMarshal(req))
		}

		err = invoker(nCtx, method, req, reply, cc, grpcOpts...)

		if opt.payload {
			if err == nil && reply != nil {
				nSpan.LogKV("Recv", opt.payloadMarshal(reply))
			}
		}

		finish(nSpan, err)
		return err
	}
}
