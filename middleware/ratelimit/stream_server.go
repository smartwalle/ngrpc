package ratelimit

import (
	"google.golang.org/grpc"
)

func WithStreamServer(limiter Limiter, opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		limiter: limiter,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerLimit(defaultOption))
}

func streamServerLimit(defaultOption *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if defaultOption.limiter != nil && defaultOption.limiter.Allow() == false {
			return errorFrom(defaultOption, info.FullMethod)
		}
		return handler(srv, ss)
	}
}
