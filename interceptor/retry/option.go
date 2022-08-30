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
	apply func(*options)
}

type options struct {
	backoff     Backoff
	codes       []codes.Code
	max         int
	callTimeout time.Duration
}

// Disable 禁用重试
func Disable() Option {
	return WithMax(0)
}

// WithMax 设置最大重试次数
func WithMax(max int) Option {
	return Option{apply: func(opts *options) {
		opts.max = max
	}}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return Option{apply: func(opts *options) {
		opts.callTimeout = timeout
	}}
}

// WithCodes 设置支持重试的错误 Code
func WithCodes(retryCodes ...codes.Code) Option {
	return Option{apply: func(opts *options) {
		opts.codes = retryCodes
	}}
}

// WithBackoff 设置重试时间延迟
func WithBackoff(f Backoff) Option {
	return Option{apply: func(opts *options) {
		opts.backoff = f
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
