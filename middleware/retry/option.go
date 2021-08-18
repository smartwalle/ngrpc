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

func Disable() Option {
	return WithMax(0)
}

func WithMax(max int) Option {
	return Option{apply: func(opt *option) {
		opt.max = max
	}}
}

func WithTimeout(timeout time.Duration) Option {
	return Option{apply: func(opt *option) {
		opt.callTimeout = timeout
	}}
}

func WithCodes(retryCodes ...codes.Code) Option {
	return Option{apply: func(opt *option) {
		opt.codes = retryCodes
	}}
}

func WithBackoff(f Backoff) Option {
	return Option{apply: func(opt *option) {
		opt.backoff = f
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

func filterOptions(inOpts []grpc.CallOption) (grpcOptions []grpc.CallOption, retryOptions []Option) {
	for _, inOpt := range inOpts {
		if opt, ok := inOpt.(Option); ok {
			retryOptions = append(retryOptions, opt)
		} else {
			grpcOptions = append(grpcOptions, inOpt)
		}
	}
	return grpcOptions, retryOptions
}
