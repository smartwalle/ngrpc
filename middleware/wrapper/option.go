package wrapper

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Handler func(ctx context.Context, md metadata.MD) (context.Context, metadata.MD)

type ErrorHandler func(err error) error

type Option struct {
	grpc.EmptyCallOption
	apply func(*options)
}

type options struct {
	handler      Handler
	errorWrapper ErrorHandler
}

func WithWrapper(h Handler) Option {
	return Option{
		apply: func(opts *options) {
			opts.handler = h
		},
	}
}

// WithErrorWrapper 错误装饰器
func WithErrorWrapper(h ErrorHandler) Option {
	return Option{
		apply: func(opts *options) {
			opts.errorWrapper = h
		},
	}
}

func defaultWrapper(ctx context.Context, md metadata.MD) (context.Context, metadata.MD) {
	return ctx, md
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

func outgoing(ctx context.Context, opts *options) context.Context {
	if opts.handler != nil {
		var md, _ = metadata.FromOutgoingContext(ctx)
		if md == nil {
			md = metadata.MD{}
		}
		ctx, md = opts.handler(ctx, md)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func incoming(ctx context.Context, opts *options) context.Context {
	if opts.handler != nil {
		var md, _ = metadata.FromIncomingContext(ctx)
		if md == nil {
			md = metadata.MD{}
		}
		ctx, md = opts.handler(ctx, md)
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	return ctx
}
