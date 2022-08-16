package hystrix

import (
	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
)

var (
	ErrMaxConcurrency = hystrix.ErrMaxConcurrency

	ErrCircuitOpen = hystrix.ErrCircuitOpen

	ErrTimeout = hystrix.ErrTimeout
)

type Filter func(err error) error

type Option struct {
	grpc.EmptyCallOption
	apply func(*options)
}

type options struct {
	filter Filter
}

// WithFilter 过滤 error
func WithFilter(h Filter) Option {
	return Option{apply: func(opts *options) {
		opts.filter = h
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
