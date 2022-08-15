package logging

import (
	"google.golang.org/grpc"
)

type Option struct {
	grpc.EmptyCallOption
	apply func(*options)
}

type options struct {
	logger  Logger
	payload bool
}

// Disable 禁用日志
func Disable() Option {
	return WithLogger(&nilLogger{})
}

// WithLogger 设置日志组件
func WithLogger(logger Logger) Option {
	if logger == nil {
		logger = &nilLogger{}
	}
	return Option{
		apply: func(opts *options) {
			opts.logger = logger
		},
	}
}

// WithPayload 设置是否需要记录请求参数及响应数据信息
func WithPayload(payload bool) Option {
	return Option{
		apply: func(opts *options) {
			opts.payload = payload
		},
	}
}

func mergeOptions(dOpts *options, callOptions []Option) *options {
	if len(callOptions) == 0 {
		return dOpts
	}
	var nOpt = &options{}
	*nOpt = *dOpts
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
