package main

import (
	"Server/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World!")
	})
	r.Run("localhost:8081")
}
