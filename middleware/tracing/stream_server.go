package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
)

func WithStreamServer(opts ...Option) grpc.ServerOption {
	var defaultOption = &option{
		tracer:         opentracing.GlobalTracer(),
		payloadMarshal: defaultPayloadMarshal,
		opName:         defaultOperationName,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerTracing(defaultOption))
}

func streamServerTracing(opt *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if opt.disable {
			return handler(srv, ss)
		}
		var opName = opt.opName(ss.Context(), info.FullMethod)

		var nCtx, err = traceSeverStreamBegin(ss.Context(), opt, opName)
		if err != nil {
			return err
		}

		var pCtx context.Context
		if opt.streamPayload {
			var pSpan opentracing.Span
			pCtx, pSpan, err = serverSpanFromContext(nCtx, opt.tracer, fmt.Sprintf("[Payload] %s", opName))
			if err != nil {
				return err
			}
			pSpan.Finish()
		}

		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
			opt:          opt,
			opName:       opName,
			pCtx:         pCtx,
		}
		err = handler(srv, nStream)

		traceServerStreamClose(nCtx, opt, opName, err)

		return err
	}
}

func traceSeverStreamBegin(ctx context.Context, opt *option, opName string) (context.Context, error) {
	var nCtx, nSpan, err = serverSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Server Stream] %s", opName))
	if err != nil {
		return nil, err
	}

	if opt.payload {
		var md, _ = metadata.FromIncomingContext(ctx)
		traceHeader(nSpan, md)
	}

	finish(nSpan, err)
	return nCtx, nil
}

func traceServerStreamClose(ctx context.Context, opt *option, opName string, err error) {
	var _, nSpan, _ = serverSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Server Stream Close] %s", opName))
	finish(nSpan, err)
}

type serverStream struct {
	grpc.ServerStream
	ctx    context.Context
	opt    *option
	opName string
	pCtx   context.Context
}

func (this *serverStream) Context() context.Context {
	return this.ctx
}

func (this *serverStream) SendMsg(m interface{}) error {
	var err = this.ServerStream.SendMsg(m)
	if err != io.EOF && this.pCtx != nil {
		traceServerStreamPayload(this.pCtx, this.opt.tracer, "Send", this.opName, m, err)
	}
	return err
}

func (this *serverStream) RecvMsg(m interface{}) error {
	var err = this.ServerStream.RecvMsg(m)
	if err != io.EOF && this.pCtx != nil {
		traceServerStreamPayload(this.pCtx, this.opt.tracer, "Recv", this.opName, m, err)
	}
	return err
}
