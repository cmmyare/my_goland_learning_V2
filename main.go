package main

import (
	"fmt"

	"github.com/cmmyare/restapi/models"
	"github.com/cmmyare/restapi/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.Use(gin.Logger())
	models.ConnectionDatabase()

  routes.RegisterRoutes(server)
	server.Run(":8080")
	fmt.Println("server started")
}