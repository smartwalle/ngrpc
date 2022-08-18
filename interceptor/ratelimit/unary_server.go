package ratelimit

import (
	"context"
	"google.golang.org/grpc"
)

// WithUnaryServer 服务端普通方法访问限流
func WithUnaryServer(limiter Limiter, opts ...Option) grpc.ServerOption {
	var defaultOptions = &options{
		limiter: limiter,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.ChainUnaryInterceptor(unaryServerLimit(defaultOptions))
}

func unaryServerLimit(opts *options) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if opts.limiter != nil && opts.limiter.Allow() == false {
			return nil, errorFrom(opts, info.FullMethod)
		}
		return handler(ctx, req)
	}
}
