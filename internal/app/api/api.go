package api

import (
	"context"
	"log/slog"
	"os"

	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/middleware"
	"github.com/Me1onRind/mr_agent/internal/initialize"
	"github.com/Me1onRind/mr_agent/internal/usecase/chat"
	"github.com/Me1onRind/mr_agent/internal/usecase/ping"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	PingUsecase *ping.PingUsecase
	ChatUsecase *chat.ChatUsecase
}

func NewAPIServer() *APIServer {
	a := &APIServer{
		PingUsecase: ping.NewPingUsecase(),
		ChatUsecase: chat.NewChatUsecase(),
	}
	return a
}

func (a *APIServer) RegisterMiddleware(r *gin.Engine) *APIServer {
	r.Use(
		middleware.Recover(),
		middleware.Tracer(),
		middleware.AccessLog(),
	)
	return a
}

func (a *APIServer) RegisterRouter(router *gin.RouterGroup) *APIServer {
	router = router.Group("/api")

	pingGroup := router.Group("/ping")
	pingGroup.POST("/echo", middleware.JSON(a.PingUsecase.Echo))
	pingGroup.POST("/panic", middleware.JSON(a.PingUsecase.Panic))
	pingGroup.POST("/hello_to_agent", middleware.JSON(a.PingUsecase.HelloToAgent))

	chatGroup := router.Group("/chat")
	chatGroup.POST("/chat", middleware.JSON(a.ChatUsecase.Chat))

	return a
}

func (a *APIServer) Init(ctx context.Context) *APIServer {
	logger.InitLogger(os.Stdout, slog.LevelDebug, true)
	_ = initialize.InitOpentracing("mr_agent", "0.0.1")(ctx)
	return a
}
