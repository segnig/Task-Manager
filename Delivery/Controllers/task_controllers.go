package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	domain "github.com/segnig/task-manager/Domains"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

func (tc *TaskController) Create(c *gin.Context) {

	var task domain.Task
	err := c.ShouldBind(&task)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	task.ID = primitive.NewObjectID()
	task.TaskID = task.ID.Hex()

	if err = tc.TaskUsecase.Create(c, &task); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Task created successfully"})
}

func (tc *TaskController) FetchAll(c *gin.Context) {

	tasks, err := tc.TaskUsecase.FetchAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (tc *TaskController) Fetch(c *gin.Context) {
	taskID := c.GetString("task_id")

	task, err := tc.TaskUsecase.FetchById(c, taskID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) Update(c *gin.Context) {
	taskID := c.GetString("task_id")

	var task domain.Task

	if err := c.ShouldBind(&task); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	if err := tc.TaskUsecase.UpdateById(c, taskID, &task); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Task updated successfully"})
}

func (tc *TaskController) Delete(c *gin.Context) {
	taskID := c.GetString("task_id")
	if err := tc.TaskUsecase.DeleteById(c, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Task Deleted successfully"})
}
