// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/recovery/interceptors.go

package recovery

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WithUnaryCall(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryRecovery(defaultOption))
}

func unaryRecovery(defaultOption *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverFrom(ctx, defaultOption, r)
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}

func recoverFrom(ctx context.Context, opt *option, r interface{}) error {
	if opt.handler != nil {
		return opt.handler(ctx, r)
	}
	return status.Errorf(codes.Internal, "%v", r)
}
