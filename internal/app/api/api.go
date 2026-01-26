package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
}

func NewApiServer() *ApiServer {
	a := &ApiServer{}
	return a
}

func (a *ApiServer) RegisterMiddleware(r *gin.Engine) *ApiServer {
	r.Use(
		middleware.Recover(),
		middleware.GinTracer(),
		middleware.AccessLog(),
	)
	return h
}

func (a *ApiServer) RegisterRouter(router *gin.RouterGroup) *ApiServer {
	return h
}

func (a *ApiServer) Init() *ApiServer {
	return h
}
