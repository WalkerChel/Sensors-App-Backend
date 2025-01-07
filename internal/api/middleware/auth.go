package middleware

import (
	"errors"
	"log"
	"net/http"
	"sensors-app/internal/api/ports"
	"sensors-app/internal/entities"
	"sensors-app/internal/service/serviceErrors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
	userIDCtx  = "userID"
)

func AuthMiddleware(env entities.Config, authService ports.Authentication) gin.HandlerFunc {
	jwtCnf := env.JWT
	return func(c *gin.Context) {
		uriPath := c.Request.URL.String()

		bearerTokenStr := c.GetHeader(authHeader)
		if bearerTokenStr == "" {
			log.Printf("AuthMiddleware no token in auth header, uri: %s", uriPath)
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
				"message": "token not found in Authorization header",
			})
			return
		}

		bearerToken := strings.Fields(bearerTokenStr)
		if len(bearerToken) != 2 {
			log.Printf("token does not consist of two separate strings, uri: %s", uriPath)
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{
				"message": "check token validity",
			})
			return
		}
		token := bearerToken[1]

		userID, err := authService.ParseToken(token, jwtCnf)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrParseToken) {
				log.Printf("AuthMiddleware err: %s, uri: %s", err, uriPath)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": err.Error(),
				})
				return
			}
			log.Printf("AuthMiddleware unknown err: %s, uri: %s", err, uriPath)
			c.AbortWithStatusJSON(http.StatusTeapot, gin.H{
				"unknown error": err.Error(),
			})
			return
		}

		equal, err := authService.CheckToken(c, userID, token)
		if err != nil {
			log.Printf("User doesn't have token to proceed request: userID: %d, uri: %s", userID, uriPath)
			if errors.Is(err, serviceErrors.ErrNoTokenForCheck) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "user is not authorized",
				})
				return
			}
			log.Printf("AuthMiddleware CheckToken err: %s, uri: %s", err, uriPath)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong when checking token",
			})
			return
		}
		if !equal {
			log.Printf("Provided token does not match the user's token in db: userID: %d, uri: %s", userID, uriPath)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "provided token does not match the user's token in db",
			})
		}

		c.Set(userIDCtx, userID)
		log.Printf("UserID was saved in request's context, userID: %d, uri: %s", userID, uriPath)

		c.Next()
	}
}
