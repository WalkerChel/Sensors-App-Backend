package serviceUser

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sensors-app/internal/entities"
	"sensors-app/internal/repository/repoErrors"
	"sensors-app/internal/service/serviceErrors"
	"sensors-app/utils"
)

const (
	salt = "v8e7545t63454"
)

type UserRepo interface {
	CreateUser(cxt context.Context, user entities.User) (int64, error)
	DeleteUser(cxt context.Context, userId int64) error
	GetUserIDByEmailAndPassword(cxt context.Context, email, password string) (int64, error)
}

type UserService struct {
	userRepo UserRepo
}

func NewUserService(userRepo UserRepo) UserService {
	return UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(cxt context.Context, user entities.User) (int64, error) {
	user.Password = utils.GeneratePasswordHash(user.Password, salt)
	id, err := s.userRepo.CreateUser(cxt, user)
	if err != nil {
		log.Printf("UserService CreateUser err: %s", err)
		if errors.Is(err, repoErrors.ErrUserAlreadyExists) {
			return 0, serviceErrors.ErrUserAlreadyExists
		}
		return 0, err
	}

	return id, nil
}

func (s *UserService) GetUserIDByEmailAndPassword(cxt context.Context, email, password string) (int64, error) {
	passwordHash := utils.GeneratePasswordHash(password, salt)

	userId, err := s.userRepo.GetUserIDByEmailAndPassword(cxt, email, passwordHash)
	if err != nil {
		if errors.Is(err, repoErrors.ErrNoUser) {
			return 0, fmt.Errorf("%w: email: %s", serviceErrors.ErrNoUserInfo, email)
		}
		return 0, err
	}

	return userId, nil
}

func (s *UserService) DeleteUser(cxt context.Context, userId int64) error {
	return s.userRepo.DeleteUser(cxt, userId)
}
