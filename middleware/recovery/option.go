package recovery

import (
	"context"
)

type Handler func(ctx context.Context, r interface{}) error

type Option struct {
	apply func(*option)
}

type option struct {
	handler Handler
}

// WithHandler 设置用于处理 panic 信息的回调函数
func WithHandler(h Handler) Option {
	return Option{apply: func(opt *option) {
		opt.handler = h
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
