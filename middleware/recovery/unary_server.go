// https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/recovery/interceptors.go

package recovery

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithUnaryServer 服务端捕获普通方法调用异常
func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerRecovery(defaultOption))
}

func unaryServerRecovery(defaultOption *option) grpc.UnaryServerInterceptor {
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
