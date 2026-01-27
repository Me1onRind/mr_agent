package api

import (
	"github.com/Me1onRind/mr_agent/internal/infrastructure/middleware"
	"github.com/Me1onRind/mr_agent/internal/usecase/ping"
	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	PingUsecase *ping.PingUsecase
}

func NewApiServer() *ApiServer {
	a := &ApiServer{
		PingUsecase: ping.NewPingUsecase(),
	}
	return a
}

func (a *ApiServer) RegisterMiddleware(r *gin.Engine) *ApiServer {
	r.Use(
		middleware.Recover(),
		middleware.Tracer(),
		middleware.AccessLog(),
	)
	return a
}

func (a *ApiServer) RegisterRouter(router *gin.RouterGroup) *ApiServer {
	router = router.Group("/api")

	pingGroup := router.Group("/ping")
	pingGroup.POST("/echo", middleware.JSON(a.PingUsecase.Echo))

	return a
}

func (a *ApiServer) Init() *ApiServer {
	return a
}
