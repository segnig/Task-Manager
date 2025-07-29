package repositories_test

import (
	"context"
	"testing"
	"time"

	domain "github.com/segnig/task-manager/Domains"
	"github.com/segnig/task-manager/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDBName     = "testdb"
	testCollection = "tasks"
	testMongoURI   = "mongodb://localhost:27017"
)

type TaskRepositoryTestSuite struct {
	suite.Suite
	db      *mongo.Database
	cleanup func()
	repo    domain.TaskRepository
}

func (suite *TaskRepositoryTestSuite) SetupSuite() {
	db, cleanup := setupTestDB(suite.T())
	suite.db = db
	suite.cleanup = cleanup
	suite.repo = repositories.NewTaskRepository(*db, testCollection)
}

func (suite *TaskRepositoryTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *TaskRepositoryTestSuite) SetupTest() {
	_, err := suite.db.Collection(testCollection).DeleteMany(context.Background(), bson.M{})
	require.NoError(suite.T(), err)
}

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(testMongoURI))
	require.NoError(t, err)

	db := client.Database(testDBName)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = db.Drop(ctx)
		_ = client.Disconnect(ctx)
	}

	return db, cleanup
}

func (suite *TaskRepositoryTestSuite) TestFetchById() {
	testTask := &domain.Task{
		TaskID:    "task1",
		Title:     "Test Task",
		CreatedBy: "user1",
	}
	_, err := suite.db.Collection(testCollection).InsertOne(context.Background(), testTask)
	require.NoError(suite.T(), err)

	suite.Run("existing task", func() {
		task, err := suite.repo.FetchById(context.Background(), "task1")
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "Test Task", task.Title)
	})

	suite.Run("non-existent task", func() {
		task, err := suite.repo.FetchById(context.Background(), "nonexistent")
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), task)
	})
}

func (suite *TaskRepositoryTestSuite) TestDeleteById() {
	testTask := &domain.Task{
		TaskID:    "task1",
		Title:     "Test Task",
		CreatedBy: "user1",
	}
	_, err := suite.db.Collection(testCollection).InsertOne(context.Background(), testTask)
	require.NoError(suite.T(), err)

	suite.Run("successful deletion", func() {
		err := suite.repo.DeleteById(context.Background(), "task1", "user1")
		assert.NoError(suite.T(), err)

		count, err := suite.db.Collection(testCollection).CountDocuments(context.Background(), bson.M{"task_id": "task1"})
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(0), count)
	})

	suite.Run("unauthorized deletion", func() {
		_, err := suite.db.Collection(testCollection).InsertOne(context.Background(), testTask)
		require.NoError(suite.T(), err)

		err = suite.repo.DeleteById(context.Background(), "task1", "wronguser")
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "unauthorized")
	})

	suite.Run("non-existent task", func() {
		err := suite.repo.DeleteById(context.Background(), "nonexistent", "user1")
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "no task found")
	})
}
func (suite *TaskRepositoryTestSuite) TestUpdateById() {
	testTask := &domain.Task{
		TaskID:      "task1",
		Title:       "Original Title",
		Description: "Original Content",
		CreatedBy:   "user1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := suite.db.Collection(testCollection).InsertOne(context.Background(), testTask)
	require.NoError(suite.T(), err)

	suite.Run("successful update", func() {
		updatedTask := &domain.Task{
			Title:       "Updated Title",
			Description: "Updated Content",
		}

		err := suite.repo.UpdateById(context.Background(), "task1", "user1", updatedTask)
		assert.NoError(suite.T(), err)

		var result domain.Task
		err = suite.db.Collection(testCollection).FindOne(context.Background(), bson.M{"task_id": "task1"}).Decode(&result)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "Updated Title", result.Title)
		assert.Equal(suite.T(), "Updated Content", result.Description)
		assert.Equal(suite.T(), "user1", result.CreatedBy)
	})

	suite.Run("unauthorized update", func() {
		updatedTask := &domain.Task{Title: "Should Fail"}
		err := suite.repo.UpdateById(context.Background(), "task1", "wronguser", updatedTask)
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "unauthorized")
	})

	suite.Run("non-existent task", func() {
		err := suite.repo.UpdateById(context.Background(), "nonexistent", "user1", &domain.Task{})
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "no task found")
	})
}

func (suite *TaskRepositoryTestSuite) TestCreate() {
	task := &domain.Task{
		TaskID:      "task1",
		Title:       "Test Task",
		Description: "Test Content",
		CreatedBy:   "user1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.Run("successful creation", func() {
		err := suite.repo.Create(context.Background(), task)
		assert.NoError(suite.T(), err)

		var result domain.Task
		err = suite.db.Collection(testCollection).FindOne(context.Background(), bson.M{"task_id": "task1"}).Decode(&result)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), task.Title, result.Title)
	})

	suite.Run("duplicate task ID", func() {

		err := suite.repo.Create(context.Background(), task)
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "already exists")
	})
}

func (suite *TaskRepositoryTestSuite) TestFetchAll() {
	testTasks := []interface{}{
		&domain.Task{TaskID: "task1", Title: "Task 1", CreatedBy: "user1"},
		&domain.Task{TaskID: "task2", Title: "Task 2", CreatedBy: "user2"},
	}
	_, err := suite.db.Collection(testCollection).InsertMany(context.Background(), testTasks)
	require.NoError(suite.T(), err)

	suite.Run("fetch all tasks", func() {
		tasks, err := suite.repo.FetchAll(context.Background())
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), tasks, 2)
	})

	suite.Run("empty collection", func() {
		_, err := suite.db.Collection(testCollection).DeleteMany(context.Background(), bson.M{})
		require.NoError(suite.T(), err)

		tasks, err := suite.repo.FetchAll(context.Background())
		assert.NoError(suite.T(), err)
		assert.Empty(suite.T(), tasks)
	})
}

func TestTaskRepositorySuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}
