package logger

import (
	"context"
	"io"
	"log/slog"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/tracer"
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

func CtxLoggerWithSpanId(ctx context.Context) *slog.Logger {
	lg := CtxLogger(ctx)
	spanId := tracer.GetSpanId(ctx)
	if len(spanId) > 0 {
		lg = lg.With(slog.String("span_id", spanId))
	}
	return lg
}

func CtxLogger(ctx context.Context) *slog.Logger {
	var lg *slog.Logger
	lgValue := ctx.Value(loggerKey{})
	if lgValue == nil {
		lg = slog.Default()
	} else {
		lg = lgValue.(*slog.Logger)
	}
	return lg
}

func WithLogger(ctx context.Context, args ...any) context.Context {
	lg := CtxLogger(ctx).With(args...)
	return context.WithValue(ctx, loggerKey{}, lg)
}
