package handler

import (
	"net/http"
	"sensors-app/internal/entities"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	UserHandlers UserHandlers
}

func (r *Handlers) InitRoutes(env entities.Config) http.Handler {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", r.UserHandlers.CreateUserHandler(env))
	}

	return router
}
