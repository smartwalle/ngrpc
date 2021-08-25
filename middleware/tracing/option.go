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

// Disable 禁用追踪
func Disable() Option {
	return Option{
		apply: func(opt *option) {
			opt.disable = true
		},
	}
}

// Enable 启用追踪
func Enable() Option {
	return Option{
		apply: func(opt *option) {
			opt.disable = false
		},
	}
}

// WithTracer 设置追踪组件
func WithTracer(tracer opentracing.Tracer) Option {
	return Option{
		apply: func(opt *option) {
			opt.tracer = tracer
		},
	}
}

// WithPayload 对于普通方法，用于设置是否需要记录请求头、请求参数及响应数据信息；对于流，用于设置建立流时是否需要记录请求头信息；
func WithPayload(payload bool) Option {
	return Option{
		apply: func(opt *option) {
			opt.payload = payload
		},
	}
}

// WithPayloadMarshal 设置请求参数及响应数据的序列化方式
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

// WithStreamPayload 设置是否需要记录流的发送和接收数据信息，只作用于流操作
func WithStreamPayload(payload bool) Option {
	return Option{
		apply: func(opt *option) {
			opt.streamPayload = payload
		},
	}
}

// WithOperationName 设置操作名称
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
