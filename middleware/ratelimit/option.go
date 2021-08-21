package ratelimit

import (
	"google.golang.org/grpc"
)

type Handler func(method string) error

type Option struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	limiter Limiter
	handler Handler
}

// WithError 设置触发限流时返回给客户端的错误信息
func WithError(h Handler) Option {
	if h == nil {
		return Option{}
	}
	return Option{
		apply: func(opt *option) {
			opt.handler = h
		},
	}
}

func mergeOptions(opt *option, callOptions []Option) *option {
	if len(callOptions) == 0 {
		return opt
	}
	var nOpt = &option{}
	*nOpt = *opt
	for _, f := range callOptions {
		if f.apply == nil {
			continue
		}
		f.apply(nOpt)
	}
	return nOpt
}
