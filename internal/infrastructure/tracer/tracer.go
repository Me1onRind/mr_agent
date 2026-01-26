package tracer

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

type spanKey struct{}
type traceIdKey struct{}
type spanIdKey struct{}

func WithTrace(ctx context.Context, traceId, spanId string) context.Context {
	return context.WithValue(context.WithValue(ctx, spanIdKey{}, spanId), traceIdKey{}, traceId)
}

func WithSpan(ctx context.Context, span opentracing.Span) context.Context {
	return context.WithValue(ctx, spanKey{}, span)
}

func GetTraceId(ctx context.Context) string {
	value := ctx.Value(traceIdKey{})
	traceId, _ := value.(string)
	return traceId
}

func GetSpanId(ctx context.Context) string {
	value := ctx.Value(spanIdKey{})
	spanId, _ := value.(string)
	return spanId
}
