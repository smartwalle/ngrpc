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
		var callOption = mergeOptions(defaultOption, traceOpts)
		if callOption.disable {
			return streamer(ctx, desc, cc, method, grpcOpts...)
		}

		var nCtx, nSpan, err = clientSpanFromContext(ctx, callOption.tracer, fmt.Sprintf("[GRPC Client Stream] %s", callOption.opName(ctx, method)))
		if err != nil {
			return nil, err
		}

		if callOption.payload {
			var md, _ = metadata.FromOutgoingContext(ctx)
			logHeader(nSpan, md)
		}

		stream, err := streamer(nCtx, desc, cc, method, grpcOpts...)

		finish(nSpan, err)
		if err != nil {
			return nil, err
		}

		var nStream = &clientStream{ClientStream: stream,
			finished: false,
			opt:      callOption,
			opName:   callOption.opName(nCtx, method),
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
