package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		handler: defaultWrapper,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerWrapper(defaultOption))
}

func streamServerWrapper(opt *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var nCtx = incoming(ss.Context(), opt)
		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
		}
		return handler(srv, nStream)
	}
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (this *serverStream) Context() context.Context {
	return this.ctx
}
