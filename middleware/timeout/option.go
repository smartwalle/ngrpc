package timeout

import (
	"google.golang.org/grpc"
	"time"
)

type Option struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	timeout time.Duration
}

// Disable 禁用超时
func Disable() Option {
	return WithValue(0)
}

// WithValue 设置超时时间
func WithValue(timeout time.Duration) Option {
	return Option{apply: func(opt *option) {
		opt.timeout = timeout
	}}
}

func mergeOptions(opt *option, callOptions []Option) *option {
	if len(callOptions) == 0 {
		return opt
	}
	var nOpt = &option{}
	*nOpt = *opt
	for _, f := range callOptions {
		f.apply(nOpt)
	}
	return nOpt
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
