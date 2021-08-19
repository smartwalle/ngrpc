package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"sync"
)

func WithStreamClient(opts ...Option) grpc.DialOption {
	var defaultOption = &option{
		tracer: opentracing.GlobalTracer(),
		opName: defaultOperationName,
	}
	defaultOption = mergeOptions(defaultOption, opts)
	return grpc.WithChainStreamInterceptor(streamClientTracing(defaultOption))
}

func streamClientTracing(defaultOption *option) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		var grpcOpts, traceOpts = filterOptions(opts)
		var opt = mergeOptions(defaultOption, traceOpts)
		if opt.disable {
			return streamer(ctx, desc, cc, method, grpcOpts...)
		}

		var nCtx, nSpan, err = clientSpanFromContext(ctx, opt.tracer, fmt.Sprintf("[GRPC Client Stream] %s", opt.opName(ctx, method)))
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
			pCtx, pSpan, _ = clientSpanFromContext(nCtx, opt.tracer, fmt.Sprintf("[Payload] %s", opt.opName(ctx, method)))
			pSpan.Finish()
		}

		var nStream = &clientStream{ClientStream: stream,
			finished: false,
			opt:      opt,
			opName:   opt.opName(ctx, method),
			pCtx:     pCtx,
		}

		return nStream, err
	}
}

type clientStream struct {
	grpc.ClientStream
	mu       sync.Mutex
	finished bool
	opt      *option
	opName   string
	pCtx     context.Context
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
	this.trace("Send", m, err)
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
	this.trace("Recv", m, err)
	if err != nil {
		this.finish(err)
	}
	return err
}

func (this *clientStream) trace(name string, m interface{}, err error) {
	if this.pCtx != nil && this.finished == false {
		var _, nSpan, _ = clientSpanFromContext(this.pCtx, this.opt.tracer, fmt.Sprintf("[%s] %s", name, this.opName))
		nSpan.LogKV(name, m)
		finish(nSpan, err)
	}
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
