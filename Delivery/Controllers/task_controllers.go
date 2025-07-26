package Controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	domain "github.com/segnig/task-manager/Domains"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	TaskUsecase domain.TaskUsecase
}

func (tc *TaskController) Create(c *gin.Context) {

	var task domain.Task
	err := c.BindJSON(&task)

	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}
	task.ID = primitive.NewObjectID()
	task.TaskID = task.ID.Hex()

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	userID, exists := c.Get("user_id")
	if !exists {
		log.Panic("user_id not found in context")
	}
	task.CreatedBy = userID.(string)
	task.UpdatedBy = userID.(string)

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
	taskID := c.Param("task_id")

	task, err := tc.TaskUsecase.FetchById(c, taskID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (tc *TaskController) Update(c *gin.Context) {
	taskID := c.Param("task_id")
	userID := c.GetString("user_id")

	var task domain.Task

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	task.UpdatedBy = userID
	task.UpdatedAt = time.Now()

	if err := tc.TaskUsecase.UpdateById(c, taskID, userID, &task); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Task updated successfully"})
}

func (tc *TaskController) Delete(c *gin.Context) {
	taskID := c.Param("task_id")
	userID := c.GetString("user_id")
	if err := tc.TaskUsecase.DeleteById(c, taskID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Task Deleted successfully"})
}
