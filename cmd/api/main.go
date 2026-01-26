package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	if err := r.Run(); err != nil {
		fmt.Println(err)
	}
}
