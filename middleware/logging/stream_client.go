package logging

import (
	"context"
	"google.golang.org/grpc"
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
		defaultOption.logger.Printf("GRPC 请求建立流 [%s - %s] \n", cc.Target(), method)

		var nStream, err = streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			defaultOption.logger.Printf("GRPC 建立到 [%s - %s] 的流发生错误 [%v] \n", cc.Target(), method, err)
		} else {
			defaultOption.logger.Printf("GRPC 建立到 [%s - %s] 的流成功 \n", cc.Target(), method)
		}
		return nStream, err
	}
}
