package logging

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// WithUnaryClient 客户端普通方法调用日志处理
func WithUnaryClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		logger: &nilLogger{},
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientLog(defaultOption))
}

func unaryClientLog(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(pCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var id, nCtx = getLogUUID(pCtx)

		defaultOption.logger.Printf("[%s] GRPC 请求访问接口 [%s - %s], 请求参数 [%v] \n", id, cc.Target(), method, req)

		var start = time.Now()
		var err = invoker(nCtx, method, req, reply, cc, opts...)
		var end = time.Now()
		if err != nil {
			defaultOption.logger.Printf("[%s] GRPC 接口 [%s - %s] 返回异常, 持续时间 [%v], 返回数据 [%v], 错误信息 [%v] \n", id, cc.Target(), method, end.Sub(start), reply, err)
		} else {
			defaultOption.logger.Printf("[%s] GRPC 接口 [%s - %s] 返回成功, 持续时间 [%v], 返回数据 [%v]\n", id, cc.Target(), method, end.Sub(start), reply)
		}
		return err
	}
}
