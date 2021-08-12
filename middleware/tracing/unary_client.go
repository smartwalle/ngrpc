package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		tracer: opentracing.GlobalTracer(),
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

		var nCtx, nSpan, err = clientSpanFromContext(ctx, callOption.tracer, fmt.Sprintf("[GRPC Client] %s", method))
		if err != nil {
			return err
		}
		err = invoker(nCtx, method, req, reply, cc, grpcOpts...)
		nSpan.Finish()
		return err
	}
}
