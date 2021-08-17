package ratelimit

import (
	"context"
	"google.golang.org/grpc"
)

func WithUnaryServer(limiter Limiter, opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		limiter: limiter,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerTracing(defaultOption))
}

func unaryServerTracing(defaultOption *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if defaultOption.limiter != nil && defaultOption.limiter.Allow() == false {
			return nil, limit(defaultOption, info.FullMethod)
		}
		return handler(ctx, req)
	}
}
