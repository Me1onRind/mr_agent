package middleware

import (
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		tx := apm.DefaultTracer.StartTransaction(c.Request.URL.Path, "api")
		defer tx.End()
		ctx = apm.ContextWithTransaction(ctx, tx)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
