package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.New()
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Hello World!")
	})
	r.Run("localhost:8081")
}
