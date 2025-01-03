package handler

import (
	"net/http"
	"sensors-app/internal/entities"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	UserHandlers UserHandlers
}

func (r *Handlers) InitRoutes(env entities.Config) http.Handler {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}

	router.Use(cors.New(corsConfig))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", r.UserHandlers.CreateUserHandler())
		auth.POST(("/sign-in"), r.UserHandlers.UserAuthenticationHandler(env))
	}

	return router
}
