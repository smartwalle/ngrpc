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
		var callOption = mergeOptions(defaultOption, traceOpts)
		if callOption.disable {
			return invoker(ctx, method, req, reply, cc, grpcOpts...)
		}

		var nCtx, nSpan, err = clientSpanFromContext(ctx, callOption.tracer, fmt.Sprintf("[GRPC Client] %s", callOption.opName(ctx, method)))
		if err != nil {
			return err
		}

		if callOption.payload {
			var md, _ = metadata.FromOutgoingContext(ctx)
			logHeader(nSpan, md)

			nSpan.LogKV("Req", req)
		}

		err = invoker(nCtx, method, req, reply, cc, grpcOpts...)

		if callOption.payload {
			if err == nil && reply != nil {
				nSpan.LogKV("Recv", reply)
			}
		}

		finish(nSpan, err)
		return err
	}
}
