package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"time"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
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
				logger.Error(ctx, "GetRawData faile", slog.String("error", err.Error()))
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
			reqHeader, err := jsoniter.MarshalToString(c.Request.Header)
			if err != nil {
				logger.Warn(ctx, "MarshalToString request header failed", slog.String("error", err.Error()))
			}
			logger.Info(ctx, "http request done",
				slog.String("client_id", c.ClientIP()),
				slog.String("method", c.Request.Method),
				slog.String("proto", c.Request.Proto),
				slog.String("host", c.Request.Host),
				slog.String("path", c.Request.RequestURI),
				slog.String("req_header", reqHeader),
				slog.String("req_body", string(request)),
				slog.String("resp_body", lw.buff.String()),
				slog.Duration("cost", end.Sub(start)),
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
