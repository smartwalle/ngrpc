package logging

import (
	"google.golang.org/grpc"
)

type CallOption struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	logger Logger
}

func Disable() CallOption {
	return WithLogger(&nilLogger{})
}

func WithLogger(logger Logger) CallOption {
	if logger == nil {
		logger = &nilLogger{}
	}
	return CallOption{
		apply: func(opt *option) {
			opt.logger = logger
		},
	}
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
