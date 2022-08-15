package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		handler: defaultWrapper,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainStreamInterceptor(streamClientWrapper(defaultOptions))
}

func streamClientWrapper(defaultOptions *options) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)
		ctx = outgoing(ctx, opt)
		return streamer(ctx, desc, cc, method, grpcOpts...)
	}
}
