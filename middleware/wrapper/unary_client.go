package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		handler: defaultWrapper,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientWrapper(defaultOption))
}

func unaryClientWrapper(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, logOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOption, logOpts)
		ctx = outgoing(ctx, opt)
		return invoker(ctx, method, req, reply, cc, grpcOpts...)
	}
}
