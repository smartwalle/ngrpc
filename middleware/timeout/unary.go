package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func WithUnaryCall(opts ...CallOption) grpc.DialOption {
	var defaultOption = &option{
		timeout: 5 * time.Second,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientTimeout(defaultOption))
}

func unaryClientTimeout(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(pCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, retryOpts = filterOptions(opts)
		var callOption = mergeOptions(defaultOption, retryOpts)

		var err error

		var nCtx, cancel = callContext(pCtx, callOption)
		defer func() {
			if cancel != nil {
				cancel()
			}
		}()
		err = invoker(nCtx, method, req, reply, cc, grpcOpts...)

		return err
	}
}

func callContext(pCtx context.Context, callOption *option) (context.Context, context.CancelFunc) {
	var nCtx = pCtx
	var cancel context.CancelFunc
	if callOption.timeout > 0 {
		nCtx, cancel = context.WithTimeout(nCtx, callOption.timeout)
	}
	return nCtx, cancel
}
