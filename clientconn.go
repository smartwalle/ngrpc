package grpc4go

import (
	"context"
	"errors"
	"github.com/smartwalle/pool4go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"time"
)

var (
	ErrServerNotFound = errors.New("server not found")
)

type DialOption struct {
	grpc.EmptyDialOption
	apply func(*option)
}

type option struct {
	poolSize    int
	dialTimeout time.Duration
}

func WithPoolSize(size int) DialOption {
	return DialOption{apply: func(opt *option) {
		opt.poolSize = size
	}}
}

func WithTimeout(timeout time.Duration) DialOption {
	return DialOption{apply: func(opt *option) {
		opt.dialTimeout = timeout
	}}
}

func filterOptions(inOpts []grpc.DialOption) (grpcOptions []grpc.DialOption, dialOptions []DialOption) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(DialOption); ok {
			dialOptions = append(dialOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, dialOptions
}

func mergeOptions(opt *option, dialOptions []DialOption) *option {
	if len(dialOptions) == 0 {
		return opt
	}
	var nOpt = &option{}
	*nOpt = *opt
	for _, f := range dialOptions {
		f.apply(nOpt)
	}
	return nOpt
}

type ClientConn struct {
	pool       pool4go.Pool
	maxRetries int
}

func Dial(target string, opts ...grpc.DialOption) *ClientConn {
	var defaultOption = &option{
		poolSize:    1,
		dialTimeout: 0,
	}

	var grpcOpts, dialOpts = filterOptions(opts)
	var dialOption = mergeOptions(defaultOption, dialOpts)

	var c = &ClientConn{}
	c.pool = pool4go.New(func() (pool4go.Conn, error) {
		var ctx = context.Background()
		if dialOption.dialTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, dialOption.dialTimeout)
			defer cancel()

			grpcOpts = append(grpcOpts, grpc.WithBlock())
		}
		return grpc.DialContext(ctx, target, grpcOpts...)
	}, pool4go.WithMaxIdle(dialOption.poolSize), pool4go.WithMaxOpen(dialOption.poolSize))
	c.maxRetries = dialOption.poolSize
	return c
}

func (this *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	for i := 0; i <= this.maxRetries; i++ {
		var conn, err = this.pool.Get()
		if err != nil {
			return err
		}

		var nConn = conn.(*grpc.ClientConn)

		var state = nConn.GetState()
		if state == connectivity.TransientFailure || state == connectivity.Shutdown {
			this.pool.Release(conn)
			continue
		}

		err = nConn.Invoke(ctx, method, args, reply, opts...)
		this.pool.Put(conn)
		return err
	}
	return ErrServerNotFound
}

func (this *ClientConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	for i := 0; i <= this.maxRetries; i++ {
		var conn, err = this.pool.Get()
		if err != nil {
			return nil, err
		}

		var nConn = conn.(*grpc.ClientConn)

		var state = nConn.GetState()
		if state == connectivity.TransientFailure || state == connectivity.Shutdown {
			this.pool.Release(conn)
			continue
		}

		stream, err := nConn.NewStream(ctx, desc, method, opts...)
		this.pool.Put(conn)
		return stream, err
	}
	return nil, ErrServerNotFound
}

func (this *ClientConn) Close() {
	this.pool.Close()
}
