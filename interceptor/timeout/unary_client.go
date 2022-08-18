package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// WithUnaryClient 客户端普通方法调用超时处理
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		timeout: 5 * time.Second,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientTimeout(defaultOptions))
}

func unaryClientTimeout(defaultOptions *options) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)

		var nCtx, cancel = callContext(ctx, opt)
		defer func() {
			if cancel != nil {
				cancel()
			}
		}()

		return invoker(nCtx, method, req, reply, cc, grpcOpts...)
	}
}

func callContext(ctx context.Context, opts *options) (context.Context, context.CancelFunc) {
	var nCtx = ctx
	var cancel context.CancelFunc
	if opts.timeout > 0 {
		nCtx, cancel = context.WithTimeout(nCtx, opts.timeout)
	}
	return nCtx, cancel
}
