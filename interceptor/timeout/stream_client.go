package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// WithStreamClient 客户端建立流超时处理
func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		timeout: 5 * time.Second,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainStreamInterceptor(streamClientTimeout(defaultOptions))
}

func streamClientTimeout(defaultOptions *options) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)

		var nCtx, cancel = callContext(ctx, opt)
		defer func() {
			if cancel != nil {
				cancel()
			}
		}()

		return streamer(nCtx, desc, cc, method, grpcOpts...)
	}
}
