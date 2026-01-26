package logger

import (
	"log/slog"
	"time"
)

var (
	loc = time.FixedZone("UTC+8", 8*3600)
)

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		t := a.Value.Time()
		return slog.String(slog.TimeKey, t.In(loc).Format("2006-01-02 15:04:05"))
	}
	return a
}
