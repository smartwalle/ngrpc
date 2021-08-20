package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type PayloadMarshal func(m interface{}) interface{}
type OperationName func(ctx context.Context, method string) string

type Option struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	tracer         opentracing.Tracer
	payload        bool
	streamPayload  bool
	payloadMarshal PayloadMarshal
	disable        bool
	opName         OperationName
}

func Disable() Option {
	return Option{
		apply: func(opt *option) {
			opt.disable = true
		},
	}
}

func Enable() Option {
	return Option{
		apply: func(opt *option) {
			opt.disable = false
		},
	}
}

func WithTracer(tracer opentracing.Tracer) Option {
	return Option{
		apply: func(opt *option) {
			opt.tracer = tracer
		},
	}
}

func WithPayload(payload bool) Option {
	return Option{
		apply: func(opt *option) {
			opt.payload = payload
		},
	}
}

func WithPayloadMarshal(h PayloadMarshal) Option {
	if h == nil {
		h = defaultPayloadMarshal
	}
	return Option{
		apply: func(opt *option) {
			opt.payloadMarshal = h
		},
	}
}

func WithStreamPayload(payload bool) Option {
	return Option{
		apply: func(opt *option) {
			opt.streamPayload = payload
		},
	}
}

func WithOperationName(h OperationName) Option {
	if h == nil {
		h = defaultOperationName
	}
	return Option{
		apply: func(opt *option) {
			opt.opName = h
		},
	}
}

func defaultPayloadMarshal(m interface{}) interface{} {
	return m
}

func defaultOperationName(ctx context.Context, method string) string {
	return method
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
