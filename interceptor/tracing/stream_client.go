package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"sync"
)

// WithStreamClient 客户端流操作追踪
func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOptions = &options{
		tracer:         opentracing.GlobalTracer(),
		payloadMarshal: defaultPayloadMarshal,
		opName:         defaultOperationName,
	}
	defaultOptions = mergeOptions(defaultOptions, opts)
	return grpc.WithChainStreamInterceptor(streamClientTracing(defaultOptions))
}

func streamClientTracing(defaultOptions *options) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var grpcOpts, nOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOptions, nOpts)
		if opt.disable {
			return streamer(ctx, desc, cc, method, grpcOpts...)
		}
		var opName = opt.opName(ctx, method)

		var nCtx, nSpan, err = clientSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Client Stream] %s", opName))
		if err != nil {
			return nil, err
		}

		if opt.payload {
			var md, _ = metadata.FromOutgoingContext(ctx)
			traceHeader(nSpan, md)
		}

		stream, err := streamer(nCtx, desc, cc, method, grpcOpts...)

		finish(nSpan, err)
		if err != nil {
			return nil, err
		}

		var pCtx context.Context
		if opt.streamPayload {
			var pSpan opentracing.Span
			pCtx, pSpan, _ = clientSpanFromContext(nCtx, opt.tracer, fmt.Sprintf("[Payload] %s", opName))
			pSpan.Finish()
		}

		var nStream = &clientStream{
			ClientStream: stream,
			finished:     false,
			opt:          opt,
			opName:       opName,
			pCtx:         pCtx,
		}

		return nStream, err
	}
}

type clientStream struct {
	grpc.ClientStream
	pCtx     context.Context
	opt      *options
	opName   string
	mu       sync.Mutex
	finished bool
}

func (this *clientStream) Header() (metadata.MD, error) {
	var header, err = this.ClientStream.Header()
	if err != nil {
		this.finish(err)
	}
	return header, err
}

func (this *clientStream) SendMsg(m interface{}) error {
	var err = this.ClientStream.SendMsg(m)
	if this.pCtx != nil {
		traceClientStreamPayload(this.pCtx, this.opt.tracer, "Send", this.opName, this.opt.payloadMarshal(m), err)
	}
	if err != nil {
		this.finish(err)
	}
	return err
}

func (this *clientStream) CloseSend() error {
	var err = this.ClientStream.CloseSend()
	this.finish(err)
	return err
}

func (this *clientStream) RecvMsg(m interface{}) error {
	var err = this.ClientStream.RecvMsg(m)
	if this.pCtx != nil && this.finished == false {
		traceClientStreamPayload(this.pCtx, this.opt.tracer, "Recv", this.opName, this.opt.payloadMarshal(m), err)
	}
	if err != nil {
		this.finish(err)
	}
	return err
}

func (this *clientStream) finish(err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.finished == false {
		var _, nSpan, _ = clientSpanFromContext(this.Context(), this.opt.tracer, fmt.Sprintf("[GRPC Client Stream Close] %s", this.opName))
		finish(nSpan, err)
		this.finished = true
	}
}
