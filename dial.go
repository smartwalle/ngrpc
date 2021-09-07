package grpc4go

import (
	"context"
	"google.golang.org/grpc"
)

func Dial(target string, opts ...grpc.DialOption) *ClientConn {
	var defaultOption = &dialOption{
		poolSize: 1,
		timeout:  0,
	}

	var grpcOpts, dialOpts = filterDialOptions(opts)
	var dialOpt = mergeDialOptions(defaultOption, dialOpts)

	var client = &ClientConn{}
	client.retry = dialOpt.poolSize
	client.pool = NewClientPool(dialOpt.poolSize, func() (*grpc.ClientConn, error) {
		var ctx = context.Background()
		if dialOpt.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, dialOpt.timeout)
			defer cancel()

			grpcOpts = append(grpcOpts, grpc.WithBlock())
		}
		return grpc.DialContext(ctx, target, grpcOpts...)
	})
	return client
}
