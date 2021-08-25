package wrapper

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Handler func(ctx context.Context, md metadata.MD) (context.Context, metadata.MD)

type Option struct {
	grpc.EmptyCallOption
	apply func(*option)
}

type option struct {
	handler Handler
}

func WithWrapper(h Handler) Option {
	return Option{
		apply: func(opt *option) {
			opt.handler = h
		},
	}
}

func defaultWrapper(ctx context.Context, md metadata.MD) (context.Context, metadata.MD) {
	return ctx, md
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

func outgoing(ctx context.Context, opt *option) context.Context {
	if opt.handler != nil {
		var md, _ = metadata.FromOutgoingContext(ctx)
		if md == nil {
			md = metadata.MD{}
		}
		ctx, md = opt.handler(ctx, md)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func incoming(ctx context.Context, opt *option) context.Context {
	if opt.handler != nil {
		var md, _ = metadata.FromIncomingContext(ctx)
		if md == nil {
			md = metadata.MD{}
		}
		ctx, md = opt.handler(ctx, md)
		ctx = metadata.NewIncomingContext(ctx, md)
	}
	return ctx
}
