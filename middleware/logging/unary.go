package logging

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

func WithUnaryCall(opts ...CallOption) grpc.DialOption {
	var defaultOption = &option{
		logger: &nilLogger{},
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainUnaryInterceptor(unaryClientLog(defaultOption))
}

func unaryClientLog(defaultOption *option) grpc.UnaryClientInterceptor {
	return func(pCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var start = time.Now()
		var err = invoker(pCtx, method, req, reply, cc, opts...)
		var end = time.Now()
		if err != nil {
			defaultOption.logger.Printf("GRPC Unary Call - %s[%v]: Req={%v}, Reply={%v}, Error={%v} \n", method, end.Sub(start), req, reply, err)
		} else {
			defaultOption.logger.Printf("GRPC Unary Call - %s[%v]: Req={%v}, Reply={%v} \n", method, end.Sub(start), req, reply)
		}

		return err
	}
}
