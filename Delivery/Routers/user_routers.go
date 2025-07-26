package Routers

import (
	"time"

	"github.com/gin-gonic/gin"
	controller "github.com/segnig/task-manager/Delivery/Controllers"
	"github.com/segnig/task-manager/Intrastructures"
	repositories "github.com/segnig/task-manager/Repositories"
	usecases "github.com/segnig/task-manager/Usecases"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	SECRET_KEY := Intrastructures.GetFromEnv("SECRET_KEY")

	mongoDB := Intrastructures.GetFromEnv("MONGO_DB")
	database := Intrastructures.DBinstance(mongoDB)
	newUserRepository := repositories.NewUserRepository(*database, "user")
	newUserUsecase := usecases.NewUserUsecase(newUserRepository, time.Duration(10*time.Second))

	newPasswordProvider := Intrastructures.NewPasswordProvider(12)

	newUserToke := Intrastructures.NeWUserToken(SECRET_KEY)

	userController := controller.UserController{
		UserUsecase: newUserUsecase,
		Password:    newPasswordProvider,
		UserToken:   newUserToke,
	}

	protected := incomingRoutes.Group("/api")
	{
		protected.Use(Intrastructures.Authentication(newUserToke))
		protected.DELETE("/users/:user_id", userController.Delete)
		protected.PUT("/users/:user_id", userController.Update)
		protected.GET("/users/:user_id", userController.Fetch)
		protected.GET("/users", userController.FetchAll)
	}
	public := incomingRoutes.Group("/api/users")
	{
		public.POST("/register", userController.Register)
		public.POST("/login", userController.Login)
	}
}
