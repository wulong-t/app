package routes

import (
	"Server/controller"
	"Server/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {
	auth := app.Group("/user")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}
	user := app.Group("/user")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/getuser/:id", controller.GetUserByID)
		user.PATCH("/updateuser/:id", controller.UpdateUser)
		user.PATCH("/follow/:id", controller.FollowUser)
	}
}
