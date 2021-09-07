package grpc4go

import (
	"context"
	"github.com/smartwalle/pool4go"
	"google.golang.org/grpc"
)

func Dial(target string, opts ...grpc.DialOption) grpc.ClientConnInterface {
	var defaultOption = &dialOption{
		poolSize: 1,
		timeout:  0,
	}

	var grpcOpts, dialOpts = filterDialOptions(opts)
	var dialOpt = mergeDialOptions(defaultOption, dialOpts)

	var c = &ClientConn{}
	c.pool = pool4go.New(func() (pool4go.Conn, error) {
		var ctx = context.Background()
		if dialOpt.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, dialOpt.timeout)
			defer cancel()

			grpcOpts = append(grpcOpts, grpc.WithBlock())
		}
		return grpc.DialContext(ctx, target, grpcOpts...)
	}, pool4go.WithMaxIdle(dialOpt.poolSize), pool4go.WithMaxOpen(dialOpt.poolSize))
	c.maxRetries = dialOpt.poolSize
	return c
}
