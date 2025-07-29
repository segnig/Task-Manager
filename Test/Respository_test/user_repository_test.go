package repositories_test

import (
	"context"
	"testing"

	domain "github.com/segnig/task-manager/Domains"
	"github.com/segnig/task-manager/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testDBNameUser     = "testdb_users"
	testCollectionUser = "users"
	testMongoURIUser   = "mongodb://localhost:27017"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db      *mongo.Database
	cleanup func()
	repo    domain.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	db, cleanup := setupTestDB(suite.T())
	suite.db = db
	suite.cleanup = cleanup
	suite.repo = repositories.NewUserRepository(*db, testCollectionUser)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	_, err := suite.db.Collection(testCollectionUser).DeleteMany(context.Background(), bson.M{})
	require.NoError(suite.T(), err)
}

func (suite *UserRepositoryTestSuite) TestCreate() {
	adminUser := &domain.User{
		UserID:   "admin1",
		Username: "admin",
		Password: "adminpass",
		UserType: "ADMIN",
	}
	adminUser.ID = primitive.NewObjectID()

	regularUser := &domain.User{
		UserID:   "user1",
		Username: "regular",
		Password: "regularpass",
		UserType: "USER",
	}
	regularUser.ID = primitive.NewObjectID()

	suite.Run("create first user as admin", func() {
		err := suite.repo.Create(context.Background(), adminUser)
		assert.NoError(suite.T(), err)
	})

	suite.Run("create regular user after admin exists", func() {
		err := suite.repo.Create(context.Background(), regularUser)
		assert.NoError(suite.T(), err)
	})

	suite.Run("duplicate user ID", func() {
		duplicate := *adminUser
		duplicate.Username = "admin2"
		err := suite.repo.Create(context.Background(), &duplicate)
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "user id 'admin1' already exists")
	})

	suite.Run("duplicate username", func() {
		duplicate := *regularUser
		duplicate.UserID = "user2"
		err := suite.repo.Create(context.Background(), &duplicate)
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "username 'regular' already exists")
	})

	suite.Run("first user must be admin", func() {
		_, err := suite.db.Collection(testCollectionUser).DeleteMany(context.Background(), bson.M{})
		require.NoError(suite.T(), err)

		nonAdmin := *regularUser
		nonAdmin.UserID = "user3"
		err = suite.repo.Create(context.Background(), &nonAdmin)
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "only an ADMIN can be the first user")
	})
}

func (suite *UserRepositoryTestSuite) TestFetchById() {
	testUser := &domain.User{
		UserID:   "test1",
		Username: "testuser",
		UserType: "ADMIN",
	}
	testUser.ID = primitive.NewObjectID()
	err := suite.repo.Create(context.Background(), testUser)
	require.NoError(suite.T(), err)

	suite.Run("existing user", func() {
		user, err := suite.repo.FetchById(context.Background(), "test1")
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "testuser", user.Username)
	})

	suite.Run("non-existent user", func() {
		user, err := suite.repo.FetchById(context.Background(), "nonexistent")
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
	})
}

func (suite *UserRepositoryTestSuite) TestGetUserByUsername() {
	testUser := &domain.User{
		UserID:   "test2",
		Username: "username_test",
		UserType: "ADMIN",
	}
	err := suite.repo.Create(context.Background(), testUser)
	require.NoError(suite.T(), err)

	suite.Run("existing username", func() {
		user, err := suite.repo.GetUserByUsername(context.Background(), "username_test")
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "test2", user.UserID)
	})

	suite.Run("non-existent username", func() {
		user, err := suite.repo.GetUserByUsername(context.Background(), "nonexistent")
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
	})
}

func (suite *UserRepositoryTestSuite) TestUpdateById() {
	testUser := &domain.User{
		UserID:    "update_test",
		FirstName: "original",
		Username:  "original",
		UserType:  "ADMIN",
	}
	err := suite.repo.Create(context.Background(), testUser)
	require.NoError(suite.T(), err)

	suite.Run("successful update", func() {
		updatedUser := &domain.User{
			FirstName: "updated",
			UserID:    "update_test",
		}
		err := suite.repo.UpdateById(context.Background(), "update_test", updatedUser)
		assert.NoError(suite.T(), err)

		user, err := suite.repo.FetchById(context.Background(), "update_test")
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "updated", user.FirstName)
	})

	suite.Run("non-existent user", func() {
		err := suite.repo.UpdateById(context.Background(), "nonexistent", &domain.User{})
		assert.Error(suite.T(), err)
		assert.Contains(suite.T(), err.Error(), "no task found")
	})
}

func (suite *UserRepositoryTestSuite) TestDeleteById() {
	testUser := &domain.User{
		UserID:   "delete_test",
		UserType: "ADMIN",
	}
	err := suite.repo.Create(context.Background(), testUser)
	require.NoError(suite.T(), err)

	suite.Run("successful deletion", func() {
		err := suite.repo.DeleteById(context.Background(), "delete_test")
		assert.NoError(suite.T(), err)

		user, err := suite.repo.FetchById(context.Background(), "delete_test")
		assert.Error(suite.T(), err)
		assert.Nil(suite.T(), user)
	})

	suite.Run("non-existent user", func() {
		err := suite.repo.DeleteById(context.Background(), "nonexistent")
		assert.NoError(suite.T(), err)
	})
}

func (suite *UserRepositoryTestSuite) TestFetchAll() {
	testUsers := []*domain.User{
		{
			ID:       primitive.NewObjectID(),
			UserID:   "user1",
			Username: "user1",
			UserType: "ADMIN",
		},
		{
			ID:       primitive.NewObjectID(),
			UserID:   "user2",
			Username: "user2",
			UserType: "USER",
		},
	}

	for _, user := range testUsers {
		err := suite.repo.Create(context.Background(), user)
		require.NoError(suite.T(), err)
	}

	suite.Run("fetch all users", func() {
		users, err := suite.repo.FetchAll(context.Background())
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), users, 2)
	})

	suite.Run("empty collection", func() {
		_, err := suite.db.Collection(testCollectionUser).DeleteMany(context.Background(), bson.M{})
		require.NoError(suite.T(), err)

		users, err := suite.repo.FetchAll(context.Background())
		assert.NoError(suite.T(), err)
		assert.Empty(suite.T(), users)
	})
}

func (suite *UserRepositoryTestSuite) TestUpdateAllToken() {
	testUser := &domain.User{
		UserID:   "token_test",
		Username: "token_user",
		UserType: "ADMIN",
	}
	err := suite.repo.Create(context.Background(), testUser)
	require.NoError(suite.T(), err)

	suite.Run("successful token update", func() {
		err := suite.repo.UpdateAllToken(
			context.Background(),
			"new_token",
			"new_refresh_token",
			"token_test",
		)
		assert.NoError(suite.T(), err)

		user, err := suite.repo.FetchById(context.Background(), "token_test")
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "new_token", user.Token)
		assert.Equal(suite.T(), "new_refresh_token", user.RefreshToken)
	})

	suite.Run("non-existent user", func() {
		err := suite.repo.UpdateAllToken(
			context.Background(),
			"token",
			"refresh",
			"nonexistent",
		)
		assert.Error(suite.T(), err)
	})
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
