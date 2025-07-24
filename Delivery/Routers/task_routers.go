package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/segnig/task-manager/Intrastructures"
)

func TaskRouters(incomingRouters *gin.Engine) {
	publicRouter := incomingRouters.Group("/api")
	publicRouter.Use(Intrastructures.Authentication())
}
