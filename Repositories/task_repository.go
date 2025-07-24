package repositories

import (
	"context"
	"fmt"

	domain "github.com/segnig/task-manager/Domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	database   mongo.Database
	collection string
}

// Create implements domains.TaskRepository.
func (tr *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	collection := tr.database.Collection(tr.collection)

	filter := bson.M{"task_id": task.TaskID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("task ID '%s' already exists", task.TaskID)
	}
	_, err = collection.InsertOne(ctx, task)
	return err
}

// DeleteById implements domains.TaskRepository.
func (tr *taskRepository) DeleteById(ctx context.Context, taskId string) error {
	collection := tr.database.Collection(tr.collection)

	objID, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return err
	}

	filter := bson.M{"task_id": objID}
	_, err = collection.DeleteOne(ctx, filter)
	return err
}

// FetchAll implements domains.TaskRepository.
func (tr *taskRepository) FetchAll(ctx context.Context) ([]*domain.Task, error) {
	collection := tr.database.Collection(tr.collection)
	filter := bson.M{}

	var tasks []*domain.Task

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &tasks)

	if err != nil {
		return nil, err
	}
	return tasks, nil

}

// FetchById implements domains.TaskRepository.
func (tr *taskRepository) FetchById(ctx context.Context, taskId string) (*domain.Task, error) {
	collection := tr.database.Collection(tr.collection)

	objID, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"task_id": objID}
	var task *domain.Task
	err = collection.FindOne(ctx, filter).Decode(&task)

	return task, err
}

// UpdateById implements domains.TaskRepository.
func (tr *taskRepository) UpdateById(ctx context.Context, taskId string, task *domain.Task) error {
	collection := tr.database.Collection(tr.collection)
	objId, err := primitive.ObjectIDFromHex(taskId)

	if err != nil {
		return err
	}

	filterStage := bson.M{"task_id": objId}
	settingStage := bson.M{"$set": &task}

	result, err := collection.UpdateOne(ctx, filterStage, settingStage)

	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no task found with id '%s'", taskId)
	}
	return nil

}

func NewTaskRepository(db mongo.Database, collection string) domain.TaskRepository {
	return &taskRepository{
		database:   db,
		collection: collection,
	}
}
