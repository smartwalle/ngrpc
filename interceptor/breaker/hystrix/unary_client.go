package hystrix

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
)

// WithUnaryClient 客户端普通方法调用添加 hystrix 熔断器
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientBreaker(defaultOptions))
}

func unaryClientBreaker(defaultOptions *options) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)
		var err = hystrix.DoC(ctx, method, func(nCtx context.Context) error {
			var err = invoker(nCtx, method, req, reply, cc, grpcOpts...)
			if err != nil && opt.filter != nil {
				return opt.filter(err)
			}
			return err
		}, nil)
		return err
	}
}
