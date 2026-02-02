package main

import (
	"context"
	"log/slog"

	"github.com/Me1onRind/mr_agent/internal/app/cli"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
)

func main() {
	ctx := context.Background()
	cliService := cli.NewCLIService()
	err := cliService.Init(ctx).Run(ctx)
	if err != nil {
		logger.LoggerFromCtx(ctx).Error("run err", slog.String("error", err.Error()))
	}
}
