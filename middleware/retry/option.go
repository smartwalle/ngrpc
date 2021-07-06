package retry

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

type Backoff func(context.Context, int) time.Duration

type CallOption struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	max     int
	timeout time.Duration
	codes   []codes.Code
	backoff Backoff
}

func Disable() CallOption {
	return WithMax(0)
}

func WithMax(max int) CallOption {
	return CallOption{apply: func(opt *option) {
		opt.max = max
	}}
}

func WithPerTimeout(timeout time.Duration) CallOption {
	return CallOption{apply: func(opt *option) {
		opt.timeout = timeout
	}}
}

func WithCodes(retryCodes ...codes.Code) CallOption {
	return CallOption{apply: func(opt *option) {
		opt.codes = retryCodes
	}}
}

func WithBackoff(f Backoff) CallOption {
	return CallOption{apply: func(opt *option) {
		opt.backoff = f
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
