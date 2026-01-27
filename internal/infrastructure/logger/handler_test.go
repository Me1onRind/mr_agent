package logger

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestLogger_Handler_Smoke(t *testing.T) {
	logger := slog.New(&mrHandler{
		writer:      os.Stdout,
		loggerLevel: slog.LevelDebug,
	})
	logger.With(
		slog.String(TraceIdKey, "trace_00001"),
		slog.String(SpanIdkey, "span_00001")).
		With(slog.String("json", `{"name":"json"}`)).
		With(slog.Duration("duration", time.Second)).
		Info("hello, world")
}
