package timeout

import (
	"google.golang.org/grpc"
	"time"
)

type CallOption struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	timeout time.Duration
}

func Disable() CallOption {
	return WithValue(0)
}

func WithValue(timeout time.Duration) CallOption {
	return CallOption{apply: func(opt *option) {
		opt.timeout = timeout
	}}
}

func mergeOptions(opt *option, callOptions []CallOption) *option {
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

func filterOptions(inOpts []grpc.CallOption) (grpcOptions []grpc.CallOption, retryOptions []CallOption) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(CallOption); ok {
			retryOptions = append(retryOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, retryOptions
}
