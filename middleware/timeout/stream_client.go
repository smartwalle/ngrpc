package timeout

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// WithStreamClient 客户端建立流超时处理
func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		timeout: 5 * time.Second,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainStreamInterceptor(streamClientTimeout(defaultOption))
}

func streamClientTimeout(defaultOption *option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var grpcOpts, timeoutOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOption, timeoutOpts)

		var nCtx, cancel = callContext(ctx, opt)
		defer func() {
			if cancel != nil {
				cancel()
			}
		}()

		return streamer(nCtx, desc, cc, method, grpcOpts...)
	}
}
