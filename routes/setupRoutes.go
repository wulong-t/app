package routes

import (
	"Server/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(app *gin.Engine) {
	app.POST("/user/register", controller.Register)
}
