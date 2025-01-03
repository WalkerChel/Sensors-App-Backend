package serviceAuth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
	"sensors-app/internal/service/serviceErrors"

	"sensors-app/utils"
	"time"
)

const (
	salt = "v8e7545t63454"
)

type TokenRepo interface {
	StoreToken(cxt context.Context, userId int64, token string, tokenTTL time.Duration) error
	DeleteToken(cxt context.Context, userId int64) error
	TokenExists(cxt context.Context, userId int64) (bool, error)
}

type UserRepo interface {
	CreateUser(cxt context.Context, user entities.User) (int64, error)
	DeleteUser(cxt context.Context, userId int64) error
	GetUserByEmailAndPassword(cxt context.Context, email, password string) (int64, error)
}

type AuthService struct {
	tokenRepo TokenRepo
	userRepo  UserRepo
}

func NewAuthService(tokenRepo TokenRepo, userRepo UserRepo) AuthService {
	return AuthService{
		tokenRepo: tokenRepo,
		userRepo:  userRepo,
	}
}

func (s *AuthService) CreateUser(cxt context.Context, user entities.User) (int64, error) {
	user.Password = utils.GeneratePasswordHash(user.Password, salt)
	id, err := s.userRepo.CreateUser(cxt, user)
	if err != nil {
		log.Printf("AuthService CreateUser err: %s", err)
		if errors.Is(err, repoErrors.ErrUserAlreadyExists) {
			return 0, serviceErrors.ErrUserAlreadyExists
		}
		return 0, err
	}

	return id, nil
}

func (s *AuthService) GetUserByEmailAndPassword(cxt context.Context, email, password string) (int64, error) {
	passwordHash := utils.GeneratePasswordHash(password, salt)

	userId, err := s.userRepo.GetUserByEmailAndPassword(cxt, email, passwordHash)
	if err != nil {
		if errors.Is(err, repoErrors.ErrNoUser) {
			return 0, fmt.Errorf("%w: email: %s", serviceErrors.ErrNoUserInfo, email)
		}
		return 0, err
	}

	return userId, nil
}

func (s *AuthService) DeleteUser(cxt context.Context, userId int64) error {
	return s.userRepo.DeleteUser(cxt, userId)
}
