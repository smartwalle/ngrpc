package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOptions = &options{
		handler: defaultWrapper,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.ChainStreamInterceptor(streamServerWrapper(defaultOptions))
}

func streamServerWrapper(opts *options) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var nCtx = incoming(ss.Context(), opts)
		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
		}
		var err = handler(srv, nStream)
		if err != nil && opts.errorWrapper != nil {
			return opts.errorWrapper(err)
		}
		return err
	}
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (stream *serverStream) Context() context.Context {
	return stream.ctx
}
