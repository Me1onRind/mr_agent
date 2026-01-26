package logger

import (
	"context"
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
	lg := ctx.Value(loggerKey{})
	if lg == nil {
		return globalLogger
	}
	return lg.(*slog.Logger)
}

func With(ctx context.Context, args ...any) context.Context {
	logger := globalLogger.With(args...)
	ctx = context.WithValue(ctx, loggerKey{}, logger)
	return ctx
}

func Debugf(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Debug(msg, args...)
}

func Infof(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Info(msg, args...)
}

func Warnf(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Warn(msg, args...)
}

func Errorf(ctx context.Context, msg string, args ...any) {
	contextLogger(ctx).Error(msg, args...)
}
