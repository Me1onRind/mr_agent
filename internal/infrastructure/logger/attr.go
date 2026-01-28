package logger

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"
)

var (
	loc = time.FixedZone("UTC+8", 8*3600)
)

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		t := a.Value.Time()
		return slog.String(slog.TimeKey, t.In(loc).Format("2006-01-02 15:04:05.000"))
	case slog.SourceKey:
		if source, ok := a.Value.Any().(*slog.Source); ok {
			fileName := filepath.Base(source.File)
			function := extractFuncName(source.Function)
			line := source.Line
			return slog.String(slog.SourceKey, fmt.Sprintf("%s %s:%d", function, fileName, line))
		}
		return a
	default:
		return a
	}
}

func extractFuncName(s string) string {
	if s == "" {
		return ""
	}
	if idx := strings.LastIndex(s, "/"); idx >= 0 && idx < len(s)-1 {
		return s[idx+1:]
	}
	return s
}
