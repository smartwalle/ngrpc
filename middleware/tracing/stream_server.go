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

		var nCtx, nSpan, err = serverSpanFromContext(ss.Context(), defaultOption.tracer, fmt.Sprintf("[GRPC Server Stream] %s", defaultOption.opName(ss.Context(), info.FullMethod)))
		if err != nil {
			return err
		}

		if defaultOption.payload {
			var md, _ = metadata.FromIncomingContext(ss.Context())
			logHeader(nSpan, md)
		}

		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
		}
		err = handler(srv, nStream)
		finish(nSpan, err)
		return err
	}
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (this *serverStream) Context() context.Context {
	return this.ctx
}
