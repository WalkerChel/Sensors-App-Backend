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
	GetUserByEmailAndPassword(cxt context.Context, email, password string) (int64, error)
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

func (h *UserHandlers) CreateUserHandler() gin.HandlerFunc {
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

		log.Printf("UserHandlers CreateUserHandler Created user with userID: %d", userID)
		c.JSON(http.StatusCreated, gin.H{"userID": userID})
	}
}

func (h *UserHandlers) UserAuthenticationHandler(env entities.Config) gin.HandlerFunc {
	jwtCnf := env.JWT
	return func(c *gin.Context) {
		var userInput apiRequests.UserAuthentication

		if err := c.BindJSON(&userInput); err != nil {
			log.Printf("UserHandlers UserAuthenticationHandler BindJSON err: %s", err)
			c.JSON(c.Writer.Status(), gin.H{"error": api.ErrInsufficientFields.Error()})
			return
		}

		userID, err := h.userService.GetUserByEmailAndPassword(c, userInput.Email, userInput.Password)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrNoUserInfo) {
				log.Printf("UserHandlers UserAuthenticationHandler ErrNoUserInfo: %s: email: %s", err, userInput.Email)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Incorrect email or password"})
				return
			}
			log.Printf("UserHandlers UserAuthenticationHandler err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong. Can not find user"})
			return
		}

		token, err := h.userService.CreateToken(c, userID, jwtCnf)
		if err != nil {
			log.Printf("UserHandlers UserAuthenticationHandler err: %s: userID: %d", err, userID)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Can not create token"})
			return
		}

		log.Printf("UserHandlers UserAuthenticationHandler Created token for userID: %d", userID)
		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})

	}
}
