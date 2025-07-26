package domains

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const TaskCollection = "task"

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `json:"title" bson:"title" validate:"required,min=4,max=50"`
	Description string             `json:"description" bson:"description" validate:"max=100"`
	Status      string             `json:"status" bson:"status"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedBy   string             `json:"created_by" bson:"created_by"`
	UpdatedBy   string             `json:"updated_by" bson:"updated_by"`
	TaskID      string             `json:"task_id" bson:"task_id"`
}

type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	FetchAll(ctx context.Context) ([]*Task, error)
	FetchById(ctx context.Context, taskId string) (*Task, error)
	UpdateById(ctx context.Context, taskId string, userID string, task *Task) error
	DeleteById(ctx context.Context, taskId string, userId string) error
}

type TaskUsecase interface {
	Create(ctx context.Context, task *Task) error
	FetchAll(ctx context.Context) ([]*Task, error)
	FetchById(ctx context.Context, taskId string) (*Task, error)
	UpdateById(ctx context.Context, taskId string, userID string, task *Task) error
	DeleteById(ctx context.Context, taskId string, userID string) error
}
