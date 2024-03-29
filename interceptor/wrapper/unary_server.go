package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOptions = &options{
		handler: defaultWrapper,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.ChainUnaryInterceptor(unaryServerWrapper(defaultOptions))
}

func unaryServerWrapper(opts *options) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = incoming(ctx, opts)
		var result, err = handler(ctx, req)
		if err != nil && opts.errorWrapper != nil {
			return result, opts.errorWrapper(err)
		}
		return result, err
	}
}
