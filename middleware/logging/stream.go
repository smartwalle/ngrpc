package logging

import (
	"context"
	"google.golang.org/grpc"
	"io"
	"time"
)

// WithStreamClient 客户端流调用日志处理
func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		logger: &nilLogger{},
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainStreamInterceptor(streamClientLog(defaultOption))
}

func streamClientLog(defaultOption *option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		defaultOption.logger.Printf("GRPC 请求建立到 [%s - %s] 的流 \n", cc.Target(), method)

		var nStream, err = streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			defaultOption.logger.Printf("GRPC 建立到 [%s - %s] 的流发生错误 [%v] \n", cc.Target(), method, err)
		} else {
			defaultOption.logger.Printf("GRPC 建立到 [%s - %s] 的流成功 \n", cc.Target(), method)
		}
		return nStream, err
	}
}

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
