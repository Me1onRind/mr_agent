package main

import (
	"context"

	"github.com/Me1onRind/mr_agent/internal/app/cli"
)

func main() {
	ctx := context.Background()
	cliService := cli.NewCLIService()
	cliService.
		Init(ctx).
		Run(ctx)
}
