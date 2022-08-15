package tracing

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/smartwalle/ngrpc"
	"google.golang.org/grpc/metadata"
	"io"
)

func clientSpanFromContext(ctx context.Context, tracer opentracing.Tracer, opName string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	var header = ngrpc.HeaderFromOutgoing(ctx)

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := tracer.Extract(opentracing.TextMap, header); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	var nSpan = tracer.StartSpan(opName, opts...)

	if err := nSpan.Tracer().Inject(nSpan.Context(), opentracing.TextMap, header); err != nil {
		return nil, nil, err
	}

	var nCtx = opentracing.ContextWithSpan(header.Context(ctx), nSpan)
	return nCtx, nSpan, nil
}

func serverSpanFromContext(ctx context.Context, tracer opentracing.Tracer, opName string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	var header = ngrpc.HeaderFromIncoming(ctx)

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := tracer.Extract(opentracing.TextMap, header); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	var nSpan = tracer.StartSpan(opName, opts...)

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

func traceClientStreamPayload(ctx context.Context, tracer opentracing.Tracer, key, opName string, payload interface{}, err error) {
	var _, nSpan, _ = clientSpanFromContext(ctx, tracer, fmt.Sprintf("[%s] %s", key, opName))
	nSpan.LogKV(key, payload)
	finish(nSpan, err)
}

func traceServerStreamPayload(ctx context.Context, tracer opentracing.Tracer, key, opName string, payload interface{}, err error) {
	var _, nSpan, _ = serverSpanFromContext(ctx, tracer, fmt.Sprintf("[%s] %s", key, opName))
	nSpan.LogKV(key, payload)
	finish(nSpan, err)
}

func traceHeader(span opentracing.Span, md metadata.MD) {
	var fields = make([]log.Field, 0, len(md))
	for key, values := range md {
		for _, value := range values {
			fields = append(fields, log.String(key, value))
		}
	}
	span.LogFields(fields...)
}
