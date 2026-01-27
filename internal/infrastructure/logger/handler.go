package logger

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	TraceIdKey = "trace_id"
	SpanIdkey  = "span_id"
	splitByte  = byte('|')
)

type mrHandler struct {
	spanId      string
	traceId     string
	loggerLevel slog.Level
	writer      io.Writer
	attrBuf     *bytes.Buffer
}

func (m *mrHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= m.loggerLevel
}

func (m *mrHandler) Handle(_ context.Context, r slog.Record) error {
	source := r.Source()
	fileName := filepath.Base(source.File)
	function := extractFuncName(source.Function)
	line := source.Line
	var buf bytes.Buffer
	_, _ = buf.WriteString(r.Time.Format("2006-01-02 15:04:05.000"))
	_ = buf.WriteByte(splitByte)
	_, _ = buf.WriteString(m.traceId)
	_ = buf.WriteByte(splitByte)
	_, _ = buf.WriteString(m.spanId)
	_ = buf.WriteByte(splitByte)
	_, _ = buf.WriteString(fileName)
	_ = buf.WriteByte(':')
	_, _ = buf.WriteString(strconv.Itoa(line))
	_ = buf.WriteByte(splitByte)
	_, _ = buf.WriteString(function)
	_ = buf.WriteByte(splitByte)
	_, _ = buf.WriteString(r.Message)
	if m.attrBuf != nil && m.attrBuf.Len() > 0 {
		_ = buf.WriteByte(splitByte)
		if _, err := buf.WriteTo(m.writer); err != nil {
			return err
		}
		m.attrBuf.WriteByte('\n')
		if _, err := m.attrBuf.WriteTo(m.writer); err != nil {
			return err
		}
	} else {
		_ = buf.WriteByte('\n')
		if _, err := buf.WriteTo(m.writer); err != nil {
			return err
		}
	}
	return nil
}

func (m *mrHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &mrHandler{
		loggerLevel: m.loggerLevel,
		writer:      m.writer,
		traceId:     m.traceId,
		spanId:      m.spanId,
	}

	if m.attrBuf != nil {
		newHandler.attrBuf = bytes.NewBuffer(m.attrBuf.Bytes())
	} else {
		newHandler.attrBuf = &bytes.Buffer{}
	}
	beginComma := newHandler.attrBuf.Len() > 0
	endCommaIdx := len(attrs) - 1

	for idx, a := range attrs {
		switch a.Key {
		case TraceIdKey:
			newHandler.traceId = a.Value.String()
		case SpanIdkey:
			newHandler.spanId = a.Value.String()
		default:
			if idx == 0 && beginComma {
				newHandler.attrBuf.WriteByte(',')
			}
			newHandler.attrBuf.WriteString(a.Key)
			newHandler.attrBuf.WriteByte('=')
			newHandler.attrBuf.WriteString(a.Value.String())
			if idx < endCommaIdx {
				newHandler.attrBuf.WriteByte(',')
			}
		}
	}
	return newHandler
}

func (m *mrHandler) WithGroup(name string) slog.Handler {
	return m
}

func extractFuncName(s string) string {
	if s == "" {
		return ""
	}
	if idx := strings.LastIndex(s, "."); idx >= 0 && idx < len(s)-1 {
		return s[idx+1:]
	}
	return s
}
