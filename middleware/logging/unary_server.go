package logging

import (
	"context"
	"google.golang.org/grpc"
)

// WithUnaryServer 服务端普通方法响应日志处理
func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		logger:  &nilLogger{},
		payload: true,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerLog(defaultOption))
}

func unaryServerLog(opt *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		opt.logger.Printf(ctx, "GRPC 收到请求: [%s], 请求参数: [%v] \n", info.FullMethod, req)

		var resp, err = handler(ctx, req)

		if opt.payload {
			if err != nil {
				opt.logger.Printf(ctx, "GRPC 处理异常: [%s], 错误信息: [%v] \n", info.FullMethod, err)
			} else {
				opt.logger.Printf(ctx, "GRPC 处理完成: [%s], 返回数据: [%v] \n", info.FullMethod, resp)
			}
		}
		return resp, err
	}
}
