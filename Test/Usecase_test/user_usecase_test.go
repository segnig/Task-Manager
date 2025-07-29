package usecases_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/segnig/task-manager/Domains"

	"github.com/segnig/task-manager/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FetchAll(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}

func (m *MockUserRepository) FetchById(ctx context.Context, userId string) (*domain.User, error) {
	args := m.Called(ctx, userId)
	var user *domain.User
	if args.Get(0) != nil {
		user = args.Get(0).(*domain.User)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) UpdateById(ctx context.Context, userId string, user *domain.User) error {
	args := m.Called(ctx, userId, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteById(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateAllToken(ctx context.Context, signedToken, signedRefreshToken, userID string) error {
	args := m.Called(ctx, signedToken, signedRefreshToken, userID)
	return args.Error(0)
}

type UserUsecaseSuite struct {
	suite.Suite
	mockRepo *MockUserRepository
	usecase  domain.UserUsecase
}

func (suite *UserUsecaseSuite) SetupTest() {
	suite.mockRepo = new(MockUserRepository)
	suite.usecase = usecases.NewUserUsecase(suite.mockRepo, time.Second*2)
}

func (suite *UserUsecaseSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *UserUsecaseSuite) TestCreate() {
	user := &domain.User{UserID: "user1", Username: "testuser"}
	suite.mockRepo.On("Create", mock.Anything, user).Return(nil)

	err := suite.usecase.Create(context.TODO(), user)
	assert.NoError(suite.T(), err)
}

func (suite *UserUsecaseSuite) TestFetchById_NotFound() {
	suite.mockRepo.On("FetchById", mock.Anything, "123").Return(nil, errors.New("not found"))

	user, err := suite.usecase.FetchById(context.TODO(), "123")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
}

func (suite *UserUsecaseSuite) TestUpdateById() {
	user := &domain.User{Username: "updated"}
	suite.mockRepo.On("UpdateById", mock.Anything, "123", user).Return(nil)

	err := suite.usecase.UpdateById(context.TODO(), "123", user)
	assert.NoError(suite.T(), err)
}

func (suite *UserUsecaseSuite) TestDeleteById() {
	suite.mockRepo.On("DeleteById", mock.Anything, "123").Return(nil)

	err := suite.usecase.DeleteById(context.TODO(), "123")
	assert.NoError(suite.T(), err)
}

func (suite *UserUsecaseSuite) TestUpdateAllToken() {
	suite.mockRepo.On("UpdateAllToken", mock.Anything, "token", "refresh", "user1").Return(nil)

	err := suite.usecase.UpdateAllToken(context.TODO(), "token", "refresh", "user1")
	assert.NoError(suite.T(), err)
}

func (suite *UserUsecaseSuite) TestGetUserByUsername() {
	user := &domain.User{Username: "john"}
	suite.mockRepo.On("GetUserByUsername", mock.Anything, "john").Return(user, nil)

	result, err := suite.usecase.GetUserByUsername(context.TODO(), "john")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Username, result.Username)
}

func (suite *UserUsecaseSuite) TestFetchAll() {
	users := []*domain.User{
		{Username: "user1"},
		{Username: "user2"},
	}
	suite.mockRepo.On("FetchAll", mock.Anything).Return(users, nil)

	result, err := suite.usecase.FetchAll(context.TODO())
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)
}

func TestUserUsecaseSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseSuite))
}
