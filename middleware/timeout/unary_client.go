package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// WithUnaryClient 普通方法调用超时处理
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		timeout: 5 * time.Second,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientTimeout(defaultOption))
}

func unaryClientTimeout(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, timeoutOpts = filterOptions(opts)
		var callOption = mergeOptions(defaultOption, timeoutOpts)

		var nCtx, cancel = callContext(ctx, callOption)
		defer func() {
			if cancel != nil {
				cancel()
			}
		}()

		return invoker(nCtx, method, req, reply, cc, grpcOpts...)
	}
}

func callContext(ctx context.Context, callOption *option) (context.Context, context.CancelFunc) {
	var nCtx = ctx
	var cancel context.CancelFunc
	if callOption.timeout > 0 {
		nCtx, cancel = context.WithTimeout(nCtx, callOption.timeout)
	}
	return nCtx, cancel
}
