package wrapper

import (
	"context"
	"google.golang.org/grpc"
)

func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		handler: defaultWrapper,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainStreamInterceptor(streamClientWrapper(defaultOption))
}

func streamClientWrapper(defaultOption *option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOption, nOpts)
		ctx = outgoing(ctx, opt)
		return streamer(ctx, desc, cc, method, grpcOpts...)
	}
}
