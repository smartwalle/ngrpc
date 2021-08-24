package logging

import (
	"context"
	"google.golang.org/grpc"
)

// WithUnaryClient 客户端普通方法调用日志处理
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		logger:  &nilLogger{},
		payload: true,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientLog(defaultOption))
}

func unaryClientLog(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var id, nCtx = getLogId(ctx)

		var grpcOpts, logOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOption, logOpts)

		opt.logger.Printf("[%s] GRPC 调用接口: [%s - %s], 请求参数: [%v] \n", id, cc.Target(), method, req)

		var err = invoker(nCtx, method, req, reply, cc, grpcOpts...)

		if opt.payload {
			if err != nil {
				opt.logger.Printf("[%s] GRPC 调用失败: [%s - %s], 错误信息: [%v] \n", id, cc.Target(), method, err)
			} else {
				opt.logger.Printf("[%s] GRPC 调用成功: [%s - %s], 返回数据: [%v] \n", id, cc.Target(), method, reply)
			}
		}
		return err
	}
}
