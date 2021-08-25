package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		handler: defaultWrapper,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerWrapper(defaultOption))
}

func unaryServerWrapper(opt *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = incoming(ctx, opt)
		return handler(ctx, req)
	}
}
