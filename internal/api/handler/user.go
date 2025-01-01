package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sensors-app/internal/api"
	apiRequests "sensors-app/internal/api/api-requests"
	"sensors-app/internal/entities"
	"sensors-app/internal/service/serviceErrors"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CheckToken(cxt context.Context, userId int64) (bool, error)
	CreateToken(cxt context.Context, userId int64, cnf entities.JWT) (string, error)
	CreateUser(cxt context.Context, user entities.User) (int64, error)
	DeleteToken(cxt context.Context, userId int64, token string) error
	DeleteUser(cxt context.Context, userId int64) error
}

type UserHandlers struct {
	userService UserService
}

func NewUserHandlers(userService UserService) UserHandlers {
	return UserHandlers{
		userService: userService,
	}
}

func (h *UserHandlers) CreateUserHandler(env entities.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userInput apiRequests.UserRegistration

		if err := c.BindJSON(&userInput); err != nil {
			log.Printf("UserHandlers CreateUserHandler BindJSON err: %s", err)
			c.JSON(c.Writer.Status(), gin.H{"error": api.ErrInsufficientFields.Error()})
			return
		}

		userID, err := h.userService.CreateUser(c,
			entities.User{
				Name:     userInput.Name,
				Email:    userInput.Email,
				Password: userInput.Password,
			})
		if err != nil {
			log.Printf("UserHandlers CreateUserHandler CreateUser err: %s", err)
			if errors.Is(err, serviceErrors.ErrUserAlreadyExists) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong. Can not create user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"userID": userID})
	}
}
