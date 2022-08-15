package retry

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

type Backoff func(context.Context, int) time.Duration

type Option struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	max         int
	callTimeout time.Duration
	codes       []codes.Code
	backoff     Backoff
}

// Disable 禁用重试
func Disable() Option {
	return WithMax(0)
}

// WithMax 设置最大重试次数
func WithMax(max int) Option {
	return Option{apply: func(opt *option) {
		opt.max = max
	}}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return Option{apply: func(opt *option) {
		opt.callTimeout = timeout
	}}
}

// WithCodes 设置支持重试的错误 Code
func WithCodes(retryCodes ...codes.Code) Option {
	return Option{apply: func(opt *option) {
		opt.codes = retryCodes
	}}
}

// WithBackoff 设置重试时间延迟
func WithBackoff(f Backoff) Option {
	return Option{apply: func(opt *option) {
		opt.backoff = f
	}}
}

func mergeOptions(dOpts *option, opts []Option) *option {
	if len(opts) == 0 {
		return dOpts
	}
	var nOpts = &option{}
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
