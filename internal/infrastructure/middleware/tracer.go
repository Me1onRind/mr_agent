package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/tracer"
	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmhttp"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceId, spanId, span := newSpanInfo(ctx, c.Request.URL.Path, c.Request.Header)
		defer span.Finish()

		ctx = tracer.WithTrace(ctx, traceId, spanId)
		ctx = tracer.WithSpan(ctx, span)
		ctx = logger.With(ctx, slog.String("trace_id", traceId))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func newSpanInfo(ctx context.Context, spanName string, header http.Header) (string, string, opentracing.Span) {
	carrier := opentracing.HTTPHeadersCarrier(header)

	var traceId, spanId string
	var span opentracing.Span
	spanContext, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)

	if spanContext == nil {
		if err != nil && !errors.Is(err, opentracing.ErrSpanContextNotFound) {
			logger.Warn(ctx, "extract span failed", slog.String("error", err.Error()), slog.Any("carrier", carrier))
		}
		span = opentracing.GlobalTracer().StartSpan(spanName)
		carrier := opentracing.HTTPHeadersCarrier{}
		if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, &carrier); err != nil {
			logger.Warn(ctx, "inject span failed", slog.String("error", err.Error()), slog.Any("carrier", carrier))
		} else {
			traceId, spanId = traceIdAndSpanIdFromSpan(carrier)
		}
	} else {
		span = opentracing.GlobalTracer().StartSpan(spanName, opentracing.ChildOf(spanContext))
		traceId, spanId = traceIdAndSpanIdFromSpan(carrier)
	}

	if len(traceId) == 0 || len(spanId) == 0 {
		logger.Warn(ctx, "one of trace_id, span_id is empty",
			slog.String("trace_id", traceId), slog.String("span_id", spanId), slog.Any("carrier", carrier))
	}

	return traceId, spanId, span
}

func traceIdAndSpanIdFromSpan(carrier opentracing.HTTPHeadersCarrier) (string, string) {
	if v, ok := carrier[apmhttp.W3CTraceparentHeader]; ok {
		return traceIdAndSpanIdFromW3CTraceparent(v)
	}
	return "", ""
}

func traceIdAndSpanIdFromW3CTraceparent(values []string) (string, string) {
	if len(values) >= 1 {
		arr := strings.Split(values[0], "-")
		if len(arr) >= 3 {
			return arr[1], arr[2]
		}
	}
	return "", ""
}
