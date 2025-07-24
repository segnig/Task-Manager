package domains

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const UserCollection = "user"

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    string             `json:"first_name" validate:"required,min=3,max=50"`
	LastName     string             `json:"last_name" validate:"required,min=3,max=50"`
	Username     string             `json:"username" validate:"required,min=5,max=25"`
	Token        string             `json:"token"`
	Password     string             `json:"password"`
	UserType     string             `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken string             `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserID       string             `json:"user_id"`
}

type SignedDetails struct {
	Username  string
	FirstName string
	LastName  string
	Uid       string
	UserType  string
	jwt.StandardClaims
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FetchAll(ctx context.Context) ([]*User, error)
	FetchById(ctx context.Context, userId string) (*User, error)
	UpdateById(ctx context.Context, userId string, user *User) error
	DeleteById(ctx context.Context, userId string) error
}

type UserUsecase interface {
	Create(ctx context.Context, user *User) error
	FetchAll(ctx context.Context) ([]*User, error)
	FetchById(ctx context.Context, userId string) (*User, error)
	UpdateById(ctx context.Context, userId string, user *User) error
	DeleteById(ctx context.Context, userId string) error
}

type PasswordServiceProvider interface {
	HashPassword(userPassword string) string
	VerifyPassword(hashedPwd, plainPwd string) (bool, string)
}

type UserToken interface {
	GenerateAllTokens(username, firstName, LastName, userType, UserID string) (signedToken, signedRefreshToken string, err error)
	UpdateAllToken(signedToken, signedRefreshToken, UserID string)
	ValidateToken(signedToken string) (claims *SignedDetails, msg string)
}
