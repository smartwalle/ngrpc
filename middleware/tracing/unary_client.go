package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		tracer: opentracing.GlobalTracer(),
		opName: defaultOperationName,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientTracing(defaultOption))
}

func unaryClientTracing(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, traceOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOption, traceOpts)
		if opt.disable {
			return invoker(ctx, method, req, reply, cc, grpcOpts...)
		}

		var nCtx, nSpan, err = clientSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Client] %s", opt.opName(ctx, method)))
		if err != nil {
			return err
		}

		if opt.payload {
			var md, _ = metadata.FromOutgoingContext(ctx)
			traceHeader(nSpan, md)

			nSpan.LogKV("Req", req)
		}

		err = invoker(nCtx, method, req, reply, cc, grpcOpts...)

		if opt.payload {
			if err == nil && reply != nil {
				nSpan.LogKV("Recv", reply)
			}
		}

		finish(nSpan, err)
		return err
	}
}
