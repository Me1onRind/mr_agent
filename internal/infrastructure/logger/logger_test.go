package logger

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type logPrint func(ctx context.Context, msg string, args ...any)

func initFixedTimeLogger(buf *bytes.Buffer) {
	globalLogger = slog.New(
		slog.NewTextHandler(
			buf,
			&slog.HandlerOptions{
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.String(slog.TimeKey, time.Unix(1769436530, 0).In(loc).Format("2006-01-02 15:04:05"))
					}
					return a
				},
				Level: slog.LevelDebug,
			},
		),
	)
}

func TestLogger_Print(t *testing.T) {
	testCauses := []struct {
		printAction logPrint
		logLevel    string
		output      string
	}{
		{
			printAction: Debugf,
			logLevel:    "DEBUG",
			output:      "time=\"2026-01-26 22:08:50\" level=DEBUG msg=test key=value\n",
		},
		{
			printAction: Infof,
			logLevel:    "INFO",
			output:      "time=\"2026-01-26 22:08:50\" level=INFO msg=test key=value\n",
		},
		{
			printAction: Warnf,
			logLevel:    "WARN",
			output:      "time=\"2026-01-26 22:08:50\" level=WARN msg=test key=value\n",
		},
		{
			printAction: Errorf,
			logLevel:    "ERROR",
			output:      "time=\"2026-01-26 22:08:50\" level=ERROR msg=test key=value\n",
		},
	}

	for _, testCase := range testCauses {
		t.Run(fmt.Sprintf("case_%s", strings.ToLower(testCase.logLevel)), func(t *testing.T) {
			ctx := context.TODO()
			var buf bytes.Buffer
			initFixedTimeLogger(&buf)
			testCase.printAction(ctx, "test", slog.String("key", "value"))
			assert.Equal(t, testCase.output, buf.String())
			//t.Log(buf.String())
		})
	}
}

func TestLogger_With(t *testing.T) {
	testCauses := []struct {
		printAction logPrint
		logLevel    string
		output      string
		with        slog.Attr
	}{
		{
			printAction: Debugf,
			logLevel:    "DEBUG",
			output:      "time=\"2026-01-26 22:08:50\" level=DEBUG msg=test request_id=test10086 key=value\n",
			with:        slog.String("request_id", "test10086"),
		},
		{
			printAction: Infof,
			logLevel:    "INFO",
			output:      "time=\"2026-01-26 22:08:50\" level=INFO msg=test count=123 key=value\n",
			with:        slog.Int("count", 123),
		},
		{
			printAction: Warnf,
			logLevel:    "WARN",
			output:      "time=\"2026-01-26 22:08:50\" level=WARN msg=test result=true key=value\n",
			with:        slog.Bool("result", true),
		},
		{
			printAction: Errorf,
			logLevel:    "ERROR",
			output:      "time=\"2026-01-26 22:08:50\" level=ERROR msg=test duration=1s key=value\n",
			with:        slog.Duration("duration", time.Second),
		},
	}

	for _, testCase := range testCauses {
		t.Run(fmt.Sprintf("case_with_%s", strings.ToLower(testCase.logLevel)), func(t *testing.T) {
			ctx := context.TODO()
			var buf bytes.Buffer
			initFixedTimeLogger(&buf)
			ctx = With(ctx, testCase.with)
			testCase.printAction(ctx, "test", slog.String("key", "value"))
			assert.Equal(t, testCase.output, buf.String())
			//t.Log(buf.String())
		})
	}
}
