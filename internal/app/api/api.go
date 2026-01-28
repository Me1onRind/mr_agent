package api

import (
	"context"
	"log/slog"
	"os"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/middleware"
	"github.com/Me1onRind/mr_agent/internal/initialize"
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
	logger.InitLogger(os.Stdout, slog.LevelDebug, true)
	initialize.InitOpentracing("mr_agent", "0.0.1")(context.TODO())
	return a
}
