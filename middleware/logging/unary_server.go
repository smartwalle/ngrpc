package logging

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// WithUnaryServer 服务端普通方法调用日志处理
func WithUnaryServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		logger: &nilLogger{},
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainUnaryInterceptor(unaryServerLog(defaultOption))
}

func unaryServerLog(defaultOption *option) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var id, nCtx = getLogUUID(ctx)

		defaultOption.logger.Printf("[%s] GRPC 接口 [%s] 收到访问请求，请求参数 [%v] \n", id, info.FullMethod, req)

		var start = time.Now()
		var resp, err = handler(nCtx, req)
		var end = time.Now()

		if err != nil {
			defaultOption.logger.Printf("[%s] GRPC 接口 [%s] 处理异常, 持续时间 [%v], 返回数据 [%v], 错误信息 [%v] \n", id, info.FullMethod, end.Sub(start), resp, err)
		} else {
			defaultOption.logger.Printf("[%s] GRPC 接口 [%s] 处理完成, 持续时间 [%v], 返回数据 [%v] \n", id, info.FullMethod, end.Sub(start), resp)
		}
		return resp, err
	}
}
