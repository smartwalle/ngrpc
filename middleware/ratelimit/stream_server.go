package ratelimit

import (
	"google.golang.org/grpc"
)

// WithStreamServer 服务端建立流限流
func WithStreamServer(limiter Limiter, opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		limiter: limiter,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerLimit(defaultOption))
}

func streamServerLimit(opt *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if opt.limiter != nil && opt.limiter.Allow() == false {
			return errorFrom(opt, info.FullMethod)
		}
		return handler(srv, ss)
	}
}
