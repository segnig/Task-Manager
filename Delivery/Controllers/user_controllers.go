package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/segnig/task-manager/Domains"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserUsecase domain.UserUsecase
}

func (uc *UserController) Create(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	userID := primitive.NewObjectID()
	user.ID = userID
	user.UserID = userID.Hex()

	if err := uc.UserUsecase.Create(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.ErrorResponse{Message: "user created successfully"})
}

func (uc *UserController) FetchAll(c *gin.Context) {
	users, err := uc.UserUsecase.FetchAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) Fetch(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := uc.UserUsecase.FetchById(c, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) Update(c *gin.Context) {
	userID := c.GetString("user_if")
	var user domain.User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	if err := uc.UserUsecase.UpdateById(c, userID, &user); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "User Update successfully"})
}

func (uc *UserController) Delete(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := uc.UserUsecase.DeleteById(c, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "user deleted successfully"})
}
