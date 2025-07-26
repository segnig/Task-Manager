package usecases

import (
	"context"
	"time"

	domain "github.com/segnig/task-manager/Domains"
)

type taskUsecase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}

// Create implements domains.TaskUsecase.
func (t *taskUsecase) Create(ctx context.Context, task *domain.Task) error {
	c, cancel := context.WithTimeout(context.Background(), t.contextTimeout)
	defer cancel()
	return t.taskRepository.Create(c, task)
}

// DeleteById implements domains.TaskUsecase.
func (t *taskUsecase) DeleteById(ctx context.Context, taskId string, userID string) error {
	c, cancel := context.WithTimeout(context.Background(), t.contextTimeout)
	defer cancel()
	return t.taskRepository.DeleteById(c, taskId, userID)
}

// FetchAll implements domains.TaskUsecase.
func (t *taskUsecase) FetchAll(ctx context.Context) ([]*domain.Task, error) {
	c, cancel := context.WithTimeout(context.Background(), t.contextTimeout)
	defer cancel()
	return t.taskRepository.FetchAll(c)
}

// FetchById implements domains.TaskUsecase.
func (t *taskUsecase) FetchById(ctx context.Context, taskId string) (*domain.Task, error) {
	c, cancel := context.WithTimeout(context.Background(), t.contextTimeout)
	defer cancel()
	return t.taskRepository.FetchById(c, taskId)
}

// UpdateById implements domains.TaskUsecase.
func (t *taskUsecase) UpdateById(ctx context.Context, taskId string, userID string, task *domain.Task) error {
	c, cancel := context.WithTimeout(context.Background(), t.contextTimeout)
	defer cancel()
	return t.taskRepository.UpdateById(c, taskId, userID, task)
}

func NewTaskUsecase(taskRepository domain.TaskRepository, contextTimeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepository: taskRepository,
		contextTimeout: contextTimeout,
	}
}
