package ratelimit

import (
	"google.golang.org/grpc"
)

type Handler func(method string) error

type Option struct {
	grpc.EmptyCallOption
	apply func(*options)
}

type options struct {
	limiter Limiter
	handler Handler
}

// WithError 设置触发限流时返回给客户端的错误信息
func WithError(h Handler) Option {
	if h == nil {
		return Option{}
	}
	return Option{
		apply: func(opts *options) {
			opts.handler = h
		},
	}
}

func mergeOptions(opts *options, callOptions []Option) *options {
	if len(callOptions) == 0 {
		return opts
	}
	var nOpt = &options{}
	*nOpt = *opts
	for _, f := range callOptions {
		if f.apply == nil {
			continue
		}
		f.apply(nOpt)
	}
	return nOpt
}
