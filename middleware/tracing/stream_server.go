package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		tracer: opentracing.GlobalTracer(),
		opName: defaultOperationName,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerTracing(defaultOption))
}

func streamServerTracing(defaultOption *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if defaultOption.disable {
			return handler(srv, ss)
		}

		var nCtx, err = streamBeginTracing(ss.Context(), defaultOption, info.FullMethod)
		if err != nil {
			return err
		}

		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
		}
		err = handler(srv, nStream)

		streamCloseTracing(nCtx, defaultOption, info.FullMethod, err)

		return err
	}
}

func streamBeginTracing(ctx context.Context, opt *option, method string) (context.Context, error) {
	var nCtx, nSpan, err = serverSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Server Stream] %s", opt.opName(ctx, method)))
	if err != nil {
		return nil, err
	}

	if opt.payload {
		var md, _ = metadata.FromIncomingContext(ctx)
		logHeader(nSpan, md)
	}

	finish(nSpan, err)
	return nCtx, nil
}

func streamCloseTracing(ctx context.Context, opt *option, method string, err error) {
	var _, nSpan, _ = serverSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Server Stream Close] %s", opt.opName(ctx, method)))
	finish(nSpan, err)
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (this *serverStream) Context() context.Context {
	return this.ctx
}
