package main

import (
	"Server/database"
	_ "Server/docs"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           Gin Golang Restspi
// @version         1.0
// @description     This is a sample server celler server based on Gin.

// @host      localhost:8081
// @BasePath  /api/v1
// @schemes   http

// @securityDefinitions.apiKey  BearerAuth
// @in header
// @name        Authorization
// @description  API Key authorization

func main() {
	database.Connect()
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Group("/api/v1")
	{
		r.GET("/", func(ctx *gin.Context) {
			ctx.String(200, "Hello World!")
		})

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run("localhost:8081")
}
