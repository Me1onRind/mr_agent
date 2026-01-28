package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request []byte
			err     error
		)

		contentType := c.ContentType()
		ctx := c.Request.Context()
		if contentType == "application/json" || contentType == "text/plain" {
			request, err = c.GetRawData()
			if err != nil {
				logger.CtxLoggerWithSpanId(ctx).Error("GetRawData failed", slog.String("error", err.Error()))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(request))
			}
		}
		lw := &logWriter{
			ResponseWriter: c.Writer,
			buff:           &bytes.Buffer{},
		}
		c.Writer = lw

		start := time.Now()
		defer func() {
			end := time.Now()
			logger.CtxLoggerWithSpanId(ctx).Info("http request done",
				slog.String("client_id", c.ClientIP()),
				slog.String("method", c.Request.Method),
				slog.String("proto", c.Request.Proto),
				slog.String("host", c.Request.Host),
				slog.String("path", c.Request.RequestURI),
				slog.Any("req_header", c.Request.Header),
				slog.String("req_body", truncateBody(string(request))),
				slog.String("resp_body", truncateBody(lw.buff.String())),
				slog.Int64("cost", end.Sub(start).Milliseconds()),
			)
		}()

		c.Next()
	}
}

type logWriter struct {
	gin.ResponseWriter
	buff *bytes.Buffer
}

func (w *logWriter) Write(b []byte) (int, error) {
	w.buff.Write(b)
	return w.ResponseWriter.Write(b)
}

func truncateBody(body string) string {
	const maxSize = 1024
	if len(body) <= maxSize {
		return body
	}
	const headSize = maxSize / 2
	const tailSize = maxSize / 2
	return body[:headSize] + "......" + body[len(body)-tailSize:]
}
