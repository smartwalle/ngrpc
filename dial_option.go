package ngrpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type DialOption struct {
	grpc.EmptyDialOption
	apply func(*dialOptions)
}

type dialOptions struct {
	poolSize int
	timeout  time.Duration
}

func WithPoolSize(size int) DialOption {
	return DialOption{apply: func(opts *dialOptions) {
		opts.poolSize = size
	}}
}

func WithTimeout(timeout time.Duration) DialOption {
	return DialOption{apply: func(opts *dialOptions) {
		opts.timeout = timeout
	}}
}

func WithInsecure() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}

func filterDialOptions(inOpts []grpc.DialOption) (grpcOptions []grpc.DialOption, dialOptions []DialOption) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(DialOption); ok {
			dialOptions = append(dialOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, dialOptions
}

func mergeDialOptions(dOpts *dialOptions, opts []DialOption) *dialOptions {
	if len(opts) == 0 {
		return dOpts
	}
	var nOpts = &dialOptions{}
	*nOpts = *dOpts
	for _, f := range opts {
		f.apply(nOpts)
	}
	return nOpts
}
