package logger

import (
	"context"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/tracer"
	"io"
	"log/slog"
	"os"
)

var globalLogger *slog.Logger = slog.New(
	slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			ReplaceAttr: replaceAttr,
			Level:       slog.LevelDebug,
		},
	),
)

type loggerKey struct{}

func InitLogger(w io.Writer, level slog.Level) {
	globalLogger = slog.New(
		slog.NewTextHandler(
			w,
			&slog.HandlerOptions{
				ReplaceAttr: replaceAttr,
				Level:       level,
			},
		),
	)
}

func contextLogger(ctx context.Context) *slog.Logger {
	lgValue := ctx.Value(loggerKey{})
	if lgValue == nil {
		lgValue = globalLogger
	}

	lg := lgValue.(*slog.Logger)

	spanId := tracer.GetSpanId(ctx)
	if len(spanId) > 0 {
		lg = lg.With(slog.String("span_id", spanId))
	}
	return lg
}

func With(ctx context.Context, args ...any) context.Context {
	logger := globalLogger.With(args...)
	ctx = context.WithValue(ctx, loggerKey{}, logger)
	return ctx
}

func Debug(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Debug(msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Info(msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Warn(msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Error(msg, args...)
}
