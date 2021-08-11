package logging

import (
	"google.golang.org/grpc"
	"io"
	"time"
)

// WithStreamServer 服务端流调用日志处理
func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		logger: &nilLogger{},
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerLog(defaultOption))
}

func streamServerLog(defaultOption *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defaultOption.logger.Printf("GRPC 流 [%s] 建立成功 \n", info.FullMethod)

		var start = time.Now()
		var err = handler(srv, ss)
		var end = time.Now()
		if err != nil && err != io.EOF {
			defaultOption.logger.Printf("GRPC 流 [%s] 异常关闭, 流持续时间 [%v], 错误信息 [%v] \n", info.FullMethod, end.Sub(start), err)
		} else {
			defaultOption.logger.Printf("GRPC 流 [%s] 正常关闭, 流持续时间 [%v] \n", info.FullMethod, end.Sub(start))
		}
		return err
	}
}
