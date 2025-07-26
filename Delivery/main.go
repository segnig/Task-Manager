package main

import (
	"github.com/gin-gonic/gin"
	routers "github.com/segnig/task-manager/Delivery/Routers"
	"github.com/segnig/task-manager/Intrastructures"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger())
	routers.TaskRoutes(router)
	routers.UserRoutes(router)

	port := Intrastructures.GetFromEnv("PORT")
	router.Run(":" + port)
}
