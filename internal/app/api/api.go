package api

import (
	"context"
	"log/slog"
	"os"

	"github.com/Me1onRind/mr_agent/internal/config"
	agentTools "github.com/Me1onRind/mr_agent/internal/infrastructure/agent/tools"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/cache/llm"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/db"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/middleware"
	"github.com/Me1onRind/mr_agent/internal/initialize"
	"github.com/Me1onRind/mr_agent/internal/usecase/chat"
	"github.com/Me1onRind/mr_agent/internal/usecase/ping"
	"github.com/Me1onRind/mr_agent/internal/usecase/tools"
	"github.com/gin-gonic/gin"
)

type APIService struct {
	PingUsecase  *ping.PingUsecase
	ChatUsecase  *chat.ChatUsecase
	ToolsUsecase *tools.ToolsUsecase
}

func NewAPIService() *APIService {
	a := &APIService{
		PingUsecase:  ping.NewPingUsecase(),
		ChatUsecase:  chat.NewChatUsecase(),
		ToolsUsecase: tools.NewToolsUsecase(),
	}
	return a
}

func (a *APIService) RegisterMiddleware(r *gin.Engine) *APIService {
	r.Use(
		middleware.Recover(),
		middleware.Tracer(),
		middleware.AccessLog(),
	)
	return a
}

func (a *APIService) RegisterRouter(router *gin.RouterGroup) *APIService {
	router = router.Group("/api")

	pingGroup := router.Group("/ping")
	pingGroup.POST("/echo", middleware.JSON(a.PingUsecase.Echo))
	pingGroup.POST("/panic", middleware.JSON(a.PingUsecase.Panic))
	pingGroup.POST("/hello_to_agent", middleware.JSON(a.PingUsecase.HelloToAgent))

	chatGroup := router.Group("/chat")
	chatGroup.POST("/chat", middleware.JSON(a.ChatUsecase.Chat))

	toolsGroup := router.Group("/tools")
	toolsGroup.POST("/call", middleware.JSON(a.ToolsUsecase.Call))

	return a
}

func (a *APIService) Init(ctx context.Context) *APIService {
	if err := config.LoadLocalConfig("./conf/config.yaml"); err != nil {
		panic(err)
	}
	logger.InitLogger(os.Stdout, slog.LevelDebug, false)
	_ = initialize.InitOpentracing("mr_agent", "0.0.1")(ctx)
	_ = llm.InitLLMCache(ctx)
	if err := db.InitRegistry(ctx, config.LocalCfg.MysqlConfigs); err != nil {
		panic(err)
	}
	if err := agentTools.InitAgentTools(ctx); err != nil {
		panic(err)
	}
	return a
}
