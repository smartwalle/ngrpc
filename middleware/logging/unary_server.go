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
		var id, nCtx = getUUID(ctx)

		opt.logger.Printf("[%s] GRPC 收到请求: [%s], 请求参数: [%v] \n", id, info.FullMethod, req)

		var resp, err = handler(nCtx, req)

		if opt.payload {
			if err != nil {
				opt.logger.Printf("[%s] GRPC 处理异常: [%s], 错误信息: [%v] \n", id, info.FullMethod, err)
			} else {
				opt.logger.Printf("[%s] GRPC 处理完成: [%s], 返回数据: [%v] \n", id, info.FullMethod, resp)
			}
		}
		return resp, err
	}
}
