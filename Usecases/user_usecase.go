package usecases

import (
	"context"
	"time"

	domain "github.com/segnig/task-manager/Domains"
)

type userUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

// Get UserByUsername implements domains.UserUsecase.
func (u *userUsecase) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	c, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.userRepository.GetUserByUsername(c, username)
}

// Create implements domains.TaskUsecase.
func (u *userUsecase) Create(ctx context.Context, user *domain.User) error {
	c, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.userRepository.Create(c, user)
}

// DeleteById implements domains.TaskUsecase.
func (u *userUsecase) DeleteById(ctx context.Context, userId string) error {
	c, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.userRepository.DeleteById(c, userId)
}

// FetchAll implements domains.TaskUsecase.
func (u *userUsecase) FetchAll(ctx context.Context) ([]*domain.User, error) {
	c, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.userRepository.FetchAll(c)
}

// FetchById implements domains.TaskUsecase.
func (u *userUsecase) FetchById(ctx context.Context, userId string) (*domain.User, error) {
	c, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.userRepository.FetchById(c, userId)
}

// UpdateById implements domains.TaskUsecase.
func (u *userUsecase) UpdateById(ctx context.Context, userId string, user *domain.User) error {
	c, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.userRepository.UpdateById(c, userId, user)
}

func (u *userUsecase) UpdateAllToken(ctx context.Context, signedToken, signedRefreshToken, UserID string) error {
	c, cancel := context.WithTimeout(context.Background(), u.contextTimeout)
	defer cancel()
	return u.userRepository.UpdateAllToken(c, signedToken, signedRefreshToken, UserID)
}

func NewUserUsecase(userRepository domain.UserRepository, contextTimeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: contextTimeout,
	}
}
