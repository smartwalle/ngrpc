package retry

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

func WithStreamCallRetry(opts ...CallOption) grpc.DialOption {
	var defaultOption = &option{
		max:         1,
		callTimeout: 5 * time.Second,
		codes:       []codes.Code{codes.ResourceExhausted, codes.Unavailable},
		backoff: func(ctx context.Context, i int) time.Duration {
			return time.Second * 1
		},
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithStreamInterceptor(streamClientRetry(defaultOption))
}

func streamClientRetry(defaultOption *option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return nil, nil
	}
}
