package logger

import (
	"context"
	"io"
	"log/slog"

	"go.elastic.co/apm"
)

type loggerKey struct{}

func InitLogger(w io.Writer, level slog.Level, json bool) {
	var lg *slog.Logger
	if json {
		lg = slog.New(
			slog.NewJSONHandler(
				w,
				&slog.HandlerOptions{
					ReplaceAttr: replaceAttr,
					Level:       level,
					AddSource:   true,
				},
			),
		)
	} else {
		lg = slog.New(
			slog.NewTextHandler(
				w,
				&slog.HandlerOptions{
					ReplaceAttr: replaceAttr,
					Level:       level,
					AddSource:   true,
				},
			),
		)
	}
	slog.SetDefault(lg)
}

func LoggerFromCtx(ctx context.Context) *slog.Logger {
	var lg *slog.Logger
	lgValue := ctx.Value(loggerKey{})
	if lgValue == nil {
		lg = slog.Default()
	} else {
		lg = lgValue.(*slog.Logger)
	}

	var traceCtx apm.TraceContext
	if span := apm.SpanFromContext(ctx); span != nil {
		traceCtx = span.TraceContext()
	} else if transaction := apm.TransactionFromContext(ctx); transaction != nil {
		traceCtx = transaction.TraceContext()
	} else {
		return lg
	}

	return lg.With(
		slog.String("trace_id", traceCtx.Trace.String()),
		slog.String("span_id", traceCtx.Span.String()),
	)
}
