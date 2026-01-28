package main

import (
	"context"
	"fmt"

	"github.com/Me1onRind/mr_agent/internal/app/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	ctx := context.Background()
	api.NewAPIServer().
		RegisterMiddleware(r).
		RegisterRouter(r.Group("/")).
		Init(ctx)

	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
