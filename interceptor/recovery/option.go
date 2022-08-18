package recovery

import (
	"context"
)

type Handler func(ctx context.Context, r interface{}) error

type Option struct {
	apply func(*options)
}

type options struct {
	handler Handler
}

// WithHandler 设置用于处理 panic 信息的回调函数
func WithHandler(h Handler) Option {
	return Option{apply: func(opts *options) {
		opts.handler = h
	}}
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
