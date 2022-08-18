package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		handler: defaultWrapper,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientWrapper(defaultOptions))
}

func unaryClientWrapper(defaultOptions *options) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)
		ctx = outgoing(ctx, opt)

		var err = invoker(ctx, method, req, reply, cc, grpcOpts...)
		if err != nil && opt.errorWrapper != nil {
			return opt.errorWrapper(err)
		}
		return err
	}
}
