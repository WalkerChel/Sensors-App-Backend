package handler

import (
	"fmt"
	"net/http"
	"sensors-app/internal/api/middleware"
	"sensors-app/internal/api/ports"
	"sensors-app/internal/entities"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	UserHandlers   UserHandlers
	RegionHandlers RegionHandlers
	SensorHandlers SensorHandlers
}

func (r *Handlers) InitRoutes(env entities.Config, authService ports.Authentication) http.Handler {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}

	router.Use(cors.New(corsConfig))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", r.UserHandlers.CreateUserHandler())
		auth.POST("/sign-in", r.UserHandlers.UserAuthenticationHandler(env, authService))
		auth.Use(middleware.AuthMiddleware(env, authService)).POST("/log-out", r.UserHandlers.UserLogOutHandler(authService))
	}

	regions := router.Group("/regions", middleware.AuthMiddleware(env, authService))
	{
		regions.GET("", r.RegionHandlers.GetAllRegionsHandler(authService))
		regions.GET(fmt.Sprintf(":%s/sensors", regionIdKey), r.SensorHandlers.GetSensorsInRegionHandler(authService))
	}

	sensors := router.Group("/sensors", middleware.AuthMiddleware(env, authService))
	{
		// /query?regionID={regionID}
		sensors.GET("", r.SensorHandlers.GetPaginatedSensorsHandler(authService))
	}

	return router
}
