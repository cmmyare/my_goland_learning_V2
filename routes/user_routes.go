package routes

import (
	"github.com/cmmyare/restapi/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(server *gin.Engine) {
	server.POST("/create_user", controllers.CreateUser)
	server.POST("/user_login", controllers.LoginUser)
}