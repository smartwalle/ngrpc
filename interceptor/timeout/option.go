package timeout

import (
	"google.golang.org/grpc"
	"time"
)

type Option struct {
	grpc.EmptyCallOption
	apply func(*options)
}

type options struct {
	timeout time.Duration
}

// Disable 禁用超时
func Disable() Option {
	return WithValue(0)
}

// WithValue 设置超时时间
func WithValue(timeout time.Duration) Option {
	return Option{apply: func(opts *options) {
		opts.timeout = timeout
	}}
}

func mergeOptions(dOpts *options, opts []Option) *options {
	if len(opts) == 0 {
		return dOpts
	}
	var nOpts = &options{}
	*nOpts = *dOpts
	for _, f := range opts {
		f.apply(nOpts)
	}
	return nOpts
}

func filterOptions(inOpts []grpc.CallOption) (grpcOptions []grpc.CallOption, nOptions []Option) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(Option); ok {
			nOptions = append(nOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, nOptions
}
