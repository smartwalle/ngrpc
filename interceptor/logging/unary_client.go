package logging

import (
	"context"
	"google.golang.org/grpc"
)

// WithUnaryClient 客户端普通方法调用日志处理
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		logger:  &nilLogger{},
		payload: true,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientLog(defaultOptions))
}

func unaryClientLog(defaultOptions *options) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)

		opt.logger.Printf(ctx, "GRPC 调用接口: [%s - %s], 请求参数: [%v] \n", cc.Target(), method, req)

		var err = invoker(ctx, method, req, reply, cc, grpcOpts...)

		if opt.payload {
			if err != nil {
				opt.logger.Printf(ctx, "GRPC 调用失败: [%s - %s], 错误信息: [%v] \n", cc.Target(), method, err)
			} else {
				opt.logger.Printf(ctx, "GRPC 调用成功: [%s - %s], 返回数据: [%v] \n", cc.Target(), method, reply)
			}
		}
		return err
	}
}
