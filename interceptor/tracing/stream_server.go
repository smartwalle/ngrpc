package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

// WithStreamServer 服务端流操作追踪
func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOptions = &options{
		tracer:         opentracing.GlobalTracer(),
		payloadMarshal: defaultPayloadMarshal,
		opName:         defaultOperationName,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.ChainStreamInterceptor(streamServerTracing(defaultOptions))
}

func streamServerTracing(opts *options) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if opts.disable {
			return handler(srv, ss)
		}
		var opName = opts.opName(ss.Context(), info.FullMethod)

		var nCtx, err = traceSeverStreamBegin(ss.Context(), opts, opName)
		if err != nil {
			return err
		}

		var pCtx context.Context
		if opts.streamPayload {
			var pSpan opentracing.Span
			pCtx, pSpan, err = serverSpanFromContext(nCtx, opts.tracer, fmt.Sprintf("[Payload] %s", opName))
			if err != nil {
				return err
			}
			pSpan.Finish()
		}

		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
			opts:         opts,
			opName:       opName,
			pCtx:         pCtx,
		}
		err = handler(srv, nStream)

		traceServerStreamClose(nCtx, opts, opName, err)

		return err
	}
}

func traceSeverStreamBegin(ctx context.Context, opts *options, opName string) (context.Context, error) {
	var nCtx, nSpan, err = serverSpanFromContext(ctx, opts.tracer, fmt.Sprintf("[GRPC Server Stream] %s", opName))
	if err != nil {
		return nil, err
	}

	if opts.payload {
		var md, _ = metadata.FromIncomingContext(ctx)
		traceHeader(nSpan, md)
	}

	finish(nSpan, err)
	return nCtx, nil
}

func traceServerStreamClose(ctx context.Context, opts *options, opName string, err error) {
	var _, nSpan, _ = serverSpanFromContext(ctx, opts.tracer, fmt.Sprintf("[GRPC Server Stream Close] %s", opName))
	finish(nSpan, err)
}

type serverStream struct {
	grpc.ServerStream
	ctx    context.Context
	pCtx   context.Context
	opts   *options
	opName string
}

func (stream *serverStream) Context() context.Context {
	return stream.ctx
}

func (stream *serverStream) SendMsg(m interface{}) error {
	var err = stream.ServerStream.SendMsg(m)
	if err != io.EOF && stream.pCtx != nil {
		traceServerStreamPayload(stream.pCtx, stream.opts.tracer, "Send", stream.opName, stream.opts.payloadMarshal(m), err)
	}
	return err
}

func (stream *serverStream) RecvMsg(m interface{}) error {
	var err = stream.ServerStream.RecvMsg(m)
	if err != io.EOF && stream.pCtx != nil {
		traceServerStreamPayload(stream.pCtx, stream.opts.tracer, "Recv", stream.opName, stream.opts.payloadMarshal(m), err)
	}
	return err
}
