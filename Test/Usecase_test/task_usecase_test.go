package usecases_test

import (
	"context"
	"testing"
	"time"

	domain "github.com/segnig/task-manager/Domains"
	"github.com/segnig/task-manager/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) FetchById(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaskRepository) UpdateById(ctx context.Context, taskId string, userID string, task *domain.Task) error {
	args := m.Called(ctx, taskId, userID, task)
	return args.Error(0)
}

func (m *MockTaskRepository) DeleteById(ctx context.Context, taskId string, userID string) error {
	args := m.Called(ctx, taskId, userID)
	return args.Error(0)
}

func (m *MockTaskRepository) FetchAll(ctx context.Context) ([]*domain.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).([]*domain.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

type TaskUsecaseSuite struct {
	suite.Suite
	mockRepo *MockTaskRepository
	usecase  domain.TaskUsecase
}

func (suite *TaskUsecaseSuite) SetupTest() {
	suite.mockRepo = new(MockTaskRepository)
	suite.usecase = usecases.NewTaskUsecase(suite.mockRepo, 2*time.Second)
}

func (suite *TaskUsecaseSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

func validTask() *domain.Task {
	return &domain.Task{
		TaskID:      "task-1",
		Title:       "Test Task",
		Description: "Test Content",
		CreatedBy:   "user-1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (suite *TaskUsecaseSuite) TestCreate_Success() {
	task := validTask()
	suite.mockRepo.On("Create", mock.Anything, task).Return(nil)

	err := suite.usecase.Create(context.Background(), task)
	assert.NoError(suite.T(), err)
}

func (suite *TaskUsecaseSuite) TestCreate_RepositoryError() {
	task := validTask()
	suite.mockRepo.On("Create", mock.Anything, task).Return(domain.ErrRepoFailure)

	err := suite.usecase.Create(context.Background(), task)
	assert.ErrorIs(suite.T(), err, domain.ErrRepoFailure)
}

func (suite *TaskUsecaseSuite) TestFetchById_Success() {
	task := validTask()
	suite.mockRepo.On("FetchById", mock.Anything, task.TaskID).Return(task, nil)

	result, err := suite.usecase.FetchById(context.Background(), task.TaskID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), task, result)
}

func (suite *TaskUsecaseSuite) TestFetchById_NotFound() {
	taskID := "non-existent-id"
	suite.mockRepo.On("FetchById", mock.Anything, taskID).Return(nil, domain.ErrTaskNotFound)

	result, err := suite.usecase.FetchById(context.Background(), taskID)
	assert.ErrorIs(suite.T(), err, domain.ErrTaskNotFound)
	assert.Nil(suite.T(), result)
}

func (suite *TaskUsecaseSuite) TestUpdateById_Success() {
	task := validTask()
	updatedTask := &domain.Task{
		Title:       "Updated Title",
		Description: "Updated Content",
	}
	suite.mockRepo.On("UpdateById", mock.Anything, task.TaskID, task.CreatedBy, updatedTask).Return(nil)

	err := suite.usecase.UpdateById(context.Background(), task.TaskID, task.CreatedBy, updatedTask)
	assert.NoError(suite.T(), err)
}

func (suite *TaskUsecaseSuite) TestUpdateById_Unauthorized() {
	task := validTask()
	updatedTask := &domain.Task{
		Title: "Updated Title",
	}
	wrongUserID := "wrong-user"
	suite.mockRepo.On("UpdateById", mock.Anything, task.TaskID, wrongUserID, updatedTask).Return(domain.ErrUnauthorized)

	err := suite.usecase.UpdateById(context.Background(), task.TaskID, wrongUserID, updatedTask)
	assert.ErrorIs(suite.T(), err, domain.ErrUnauthorized)
}

func (suite *TaskUsecaseSuite) TestDeleteById_Success() {
	task := validTask()
	suite.mockRepo.On("DeleteById", mock.Anything, task.TaskID, task.CreatedBy).Return(nil)

	err := suite.usecase.DeleteById(context.Background(), task.TaskID, task.CreatedBy)
	assert.NoError(suite.T(), err)
}

func (suite *TaskUsecaseSuite) TestDeleteById_NotFound() {
	taskID := "non-existent-id"
	userID := "user-1"
	suite.mockRepo.On("DeleteById", mock.Anything, taskID, userID).Return(domain.ErrTaskNotFound)

	err := suite.usecase.DeleteById(context.Background(), taskID, userID)
	assert.ErrorIs(suite.T(), err, domain.ErrTaskNotFound)
}

func (suite *TaskUsecaseSuite) TestFetchAll_Success() {
	tasks := []*domain.Task{
		validTask(),
		{
			TaskID:    "task-2",
			Title:     "Another Task",
			CreatedBy: "user-1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	suite.mockRepo.On("FetchAll", mock.Anything).Return(tasks, nil)

	result, err := suite.usecase.FetchAll(context.Background())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), tasks, result)
}

func (suite *TaskUsecaseSuite) TestFetchAll_Empty() {
	suite.mockRepo.On("FetchAll", mock.Anything).Return([]*domain.Task{}, nil)

	result, err := suite.usecase.FetchAll(context.Background())
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

func TestTaskUsecaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUsecaseSuite))
}
