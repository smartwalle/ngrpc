package ratelimit

import (
	"google.golang.org/grpc"
)

// WithStreamServer 服务端建立流限流
func WithStreamServer(limiter Limiter, opts ...Option) grpc.ServerOption {
	var defaultOptions = &options{
		limiter: limiter,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.ChainStreamInterceptor(streamServerLimit(defaultOptions))
}

func streamServerLimit(opts *options) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if opts.limiter != nil && opts.limiter.Allow() == false {
			return errorFrom(opts, info.FullMethod)
		}
		return handler(srv, ss)
	}
}
