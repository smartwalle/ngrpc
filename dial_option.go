package grpc4go

import (
	"google.golang.org/grpc"
	"time"
)

type DialOption struct {
	grpc.EmptyDialOption
	apply func(*dialOption)
}

type dialOption struct {
	poolSize int
	timeout  time.Duration
}

func WithPoolSize(size int) DialOption {
	return DialOption{apply: func(opt *dialOption) {
		opt.poolSize = size
	}}
}

func WithTimeout(timeout time.Duration) DialOption {
	return DialOption{apply: func(opt *dialOption) {
		opt.timeout = timeout
	}}
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

func mergeDialOptions(dialOpt *dialOption, dialOptions []DialOption) *dialOption {
	if len(dialOptions) == 0 {
		return dialOpt
	}
	var nOpt = &dialOption{}
	*nOpt = *dialOpt
	for _, f := range dialOptions {
		f.apply(nOpt)
	}
	return nOpt
}
