package ngrpc

import (
	"context"
	"google.golang.org/grpc"
)

func Dial(target string, opts ...grpc.DialOption) *ClientConn {
	var defaultOption = &dialOptions{
		poolSize: 1,
		timeout:  0,
	}

	var grpcOpts, dialOpts = filterDialOptions(opts)
	var nDialOpts = mergeDialOptions(defaultOption, dialOpts)

	var client = &ClientConn{}
	client.retry = nDialOpts.poolSize
	client.pool = NewClientPool(nDialOpts.poolSize, func() (*grpc.ClientConn, error) {
		var ctx = context.Background()
		if nDialOpts.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, nDialOpts.timeout)
			defer cancel()

			grpcOpts = append(grpcOpts, grpc.WithBlock())
		}
		return grpc.DialContext(ctx, target, grpcOpts...)
	})
	return client
}
