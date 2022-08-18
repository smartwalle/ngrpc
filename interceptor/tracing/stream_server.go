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
	opts   *options
	opName string
	pCtx   context.Context
}

func (this *serverStream) Context() context.Context {
	return this.ctx
}

func (this *serverStream) SendMsg(m interface{}) error {
	var err = this.ServerStream.SendMsg(m)
	if err != io.EOF && this.pCtx != nil {
		traceServerStreamPayload(this.pCtx, this.opts.tracer, "Send", this.opName, this.opts.payloadMarshal(m), err)
	}
	return err
}

func (this *serverStream) RecvMsg(m interface{}) error {
	var err = this.ServerStream.RecvMsg(m)
	if err != io.EOF && this.pCtx != nil {
		traceServerStreamPayload(this.pCtx, this.opts.tracer, "Recv", this.opName, this.opts.payloadMarshal(m), err)
	}
	return err
}
