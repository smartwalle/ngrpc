package tracing

import (
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type Option struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	tracer  opentracing.Tracer
	disable bool
}

func Disable() Option {
	return Option{
		apply: func(opt *option) {
			opt.disable = true
		},
	}
}

func WithTracer(tracer opentracing.Tracer) Option {
	return Option{
		apply: func(opt *option) {
			opt.tracer = tracer
			opt.disable = false
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
