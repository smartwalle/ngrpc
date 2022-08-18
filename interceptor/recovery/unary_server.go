package recovery

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithUnaryServer 服务端捕获普通方法调用异常
func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &options{}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerRecovery(defaultOption))
}

func unaryServerRecovery(opts *options) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errorFrom(ctx, opts, r)
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func errorFrom(ctx context.Context, opts *options, r interface{}) error {
	if opts.handler != nil {
		return opts.handler(ctx, r)
	}
	return status.Errorf(codes.Internal, "%v", r)
}
