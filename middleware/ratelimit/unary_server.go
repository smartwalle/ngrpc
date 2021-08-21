package ratelimit

import (
	"context"
	"google.golang.org/grpc"
)

// WithUnaryServer 服务端普通方法访问限流
func WithUnaryServer(limiter Limiter, opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		limiter: limiter,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerLimit(defaultOption))
}

func unaryServerLimit(opt *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if opt.limiter != nil && opt.limiter.Allow() == false {
			return nil, errorFrom(opt, info.FullMethod)
		}
		return handler(ctx, req)
	}
}
