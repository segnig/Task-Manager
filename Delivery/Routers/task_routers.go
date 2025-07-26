package Routers

import (
	"time"

	"github.com/gin-gonic/gin"
	controller "github.com/segnig/task-manager/Delivery/Controllers"
	"github.com/segnig/task-manager/Intrastructures"
	repositories "github.com/segnig/task-manager/Repositories"
	usecases "github.com/segnig/task-manager/Usecases"
)

func TaskRoutes(incomingRoutes *gin.Engine) {

	SECRET_KEY := Intrastructures.GetFromEnv("SECRET_KEY")
	ut := Intrastructures.NeWUserToken(SECRET_KEY)

	mongoDB := Intrastructures.GetFromEnv("MONGO_DB")
	database := Intrastructures.DBinstance(mongoDB)
	newTaskRepository := repositories.NewTaskRepository(*database, "task")
	newTaskUsecase := usecases.NewTaskUsecase(newTaskRepository, time.Duration(10*time.Second))

	taskController := controller.TaskController{TaskUsecase: newTaskUsecase}

	protected := incomingRoutes.Group("/api")
	{
		protected.Use(Intrastructures.Authentication(ut))
		protected.DELETE("/tasks/:task_id", taskController.Delete)
		protected.PUT("/tasks/:task_id", taskController.Update)
		protected.POST("/tasks", taskController.Create)
	}
	public := incomingRoutes.Group("/api")
	{
		public.GET("/tasks/:task_id", taskController.Fetch)
		public.GET("/tasks", taskController.FetchAll)
	}
}
