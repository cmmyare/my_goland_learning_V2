package routes

import (
	"github.com/cmmyare/restapi/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(server *gin.Engine) {
	server.POST("/create_task", controllers.CreateTask)
}
