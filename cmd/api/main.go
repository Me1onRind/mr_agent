package main

import (
	"fmt"

	"github.com/Me1onRind/mr_agent/internal/app/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	api.NewAPIServer().
		RegisterMiddleware(r).
		RegisterRouter(r.Group("/")).
		Init()

	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
