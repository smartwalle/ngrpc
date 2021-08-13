package tracing

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/smartwalle/grpc4go"
	"io"
)

func clientSpanFromContext(ctx context.Context, tracer opentracing.Tracer, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	var header = grpc4go.HeaderFromOutgoing(ctx)

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := tracer.Extract(opentracing.TextMap, header); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	var nSpan = tracer.StartSpan(name, opts...)

	if err := nSpan.Tracer().Inject(nSpan.Context(), opentracing.TextMap, header); err != nil {
		return nil, nil, err
	}

	var nCtx = opentracing.ContextWithSpan(header.Context(ctx), nSpan)
	return nCtx, nSpan, nil
}

func serverSpanFromContext(ctx context.Context, tracer opentracing.Tracer, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	var header = grpc4go.HeaderFromIncoming(ctx)

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := tracer.Extract(opentracing.TextMap, header); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	var nSpan = tracer.StartSpan(name, opts...)

	if err := nSpan.Tracer().Inject(nSpan.Context(), opentracing.TextMap, header); err != nil {
		return nil, nil, err
	}

	var nCtx = opentracing.ContextWithSpan(header.Context(ctx), nSpan)
	return nCtx, nSpan, nil
}

func finish(span opentracing.Span, err error) {
	if err != nil && err != io.EOF {
		ext.Error.Set(span, true)
		span.LogKV("error", err.Error())
	}
	span.Finish()
}
