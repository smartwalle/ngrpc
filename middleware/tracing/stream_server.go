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
		tracer: opentracing.GlobalTracer(),
		opName: defaultOperationName,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.ChainStreamInterceptor(streamServerTracing(defaultOption))
}

func streamServerTracing(opt *option) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if opt.disable {
			return handler(srv, ss)
		}

		var nCtx, err = traceSeverStreamBegin(ss.Context(), opt, info.FullMethod)
		if err != nil {
			return err
		}

		var pCtx context.Context
		if opt.streamPayload {
			var pSpan opentracing.Span
			pCtx, pSpan, _ = serverSpanFromContext(nCtx, opt.tracer, fmt.Sprintf("[Payload] %s", opt.opName(ss.Context(), info.FullMethod)))
			pSpan.Finish()
		}

		var nStream = &serverStream{
			ServerStream: ss,
			ctx:          nCtx,
			opt:          opt,
			opName:       opt.opName(ss.Context(), info.FullMethod),
			pCtx:         pCtx,
		}
		err = handler(srv, nStream)

		traceServerStreamClose(nCtx, opt, info.FullMethod, err)

		return err
	}
}

func traceSeverStreamBegin(ctx context.Context, opt *option, method string) (context.Context, error) {
	var nCtx, nSpan, err = serverSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Server Stream] %s", opt.opName(ctx, method)))
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

func traceServerStreamClose(ctx context.Context, opt *option, method string, err error) {
	var _, nSpan, _ = serverSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Server Stream Close] %s", opt.opName(ctx, method)))
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
	this.trace("Send", m, err)
	return err
}

func (this *serverStream) RecvMsg(m interface{}) error {
	var err = this.ServerStream.RecvMsg(m)
	this.trace("Recv", m, err)
	return err
}

func (this *serverStream) trace(name string, m interface{}, err error) {
	if err == io.EOF {
		return
	}

	if this.pCtx != nil {
		var _, nSpan, _ = clientSpanFromContext(this.pCtx, this.opt.tracer, fmt.Sprintf("[%s] %s", name, this.opName))
		nSpan.LogKV(name, m)
		finish(nSpan, err)
	}
}
