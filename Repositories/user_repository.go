package repositories

import (
	"context"
	"fmt"
	"log"

	domain "github.com/segnig/task-manager/Domains"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	database   mongo.Database
	collection string
}

func (ur *userRepository) DeleteById(ctx context.Context, userId string) error {
	collection := ur.database.Collection(ur.collection)

	// objID, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return err
	// }

	filter := bson.M{"userid": userId}
	_, err := collection.DeleteOne(ctx, filter)
	return err
}

// FetchAll implements domains.userRepository.
func (ur *userRepository) FetchAll(ctx context.Context) ([]*domain.User, error) {
	collection := ur.database.Collection(ur.collection)
	filter := bson.M{}

	var users []*domain.User

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &users)

	if err != nil {
		return nil, err
	}
	return users, nil

}

// FetchById implements domains.userRepository.
func (ur *userRepository) FetchById(ctx context.Context, userId string) (*domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	// objID, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	return nil, err
	// }

	filter := bson.M{"userid": userId}
	log.Println("UserID in code block: ", userId)
	var user *domain.User
	err := collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}

// UpdateById implements domains.userRepository.
func (ur *userRepository) UpdateById(ctx context.Context, userId string, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)
	objId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return err
	}

	filterStage := bson.M{"userid": objId}
	settingStage := bson.M{"$set": &user}

	result, err := collection.UpdateOne(ctx, filterStage, settingStage)

	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no task found with id '%s'", userId)
	}
	return nil

}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"user_id": user.UserID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("user id '%s' already exists", user.UserID)
	}
	filter = bson.M{"username": user.Username}
	count, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("username '%s' already exists", user.Username)
	}

	totalUsers, err := collection.CountDocuments(ctx, bson.D{})

	if err != nil {
		return err
	}

	if totalUsers == 0 && user.UserType != "ADMIN" {
		return fmt.Errorf("only an ADMIN can be the first user")
	}
	_, err = collection.InsertOne(ctx, user)
	return err
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"username": username}
	var user *domain.User
	err := collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (ur *userRepository) UpdateAllToken(ctx context.Context, signedToken, signedRefreshToken, UserID string) error {
	user, err := ur.FetchById(ctx, UserID)
	if err != nil {
		return err
	}

	user.Token = signedToken
	user.RefreshToken = signedRefreshToken
	err = ur.UpdateById(ctx, UserID, user)
	return err
}

func NewUserRepository(db mongo.Database, collection string) *userRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}
