package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sensors-app/internal/api/ports"
	"sensors-app/internal/entities"
	"sensors-app/internal/service/serviceErrors"

	"github.com/gin-gonic/gin"
)

type RegionService interface {
	GetAllRegions(ctx context.Context) ([]entities.Region, error)
}

type RegionHandlers struct {
	regionService RegionService
}

func NewRegionService(regionService RegionService) RegionHandlers {
	return RegionHandlers{
		regionService: regionService,
	}
}

// /regions
func (h *RegionHandlers) GetAllRegionsHandler(authService ports.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := authService.GetUserIDFromCtx(c, userIDCtxKey)
		if err != nil {
			log.Printf("RegionHandlers GetAllRegionsHandler err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong when extracting userID from request's context",
			})
			return
		}

		regions, err := h.regionService.GetAllRegions(c)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrNoRegionsData) {
				log.Printf("RegionHandlers GetAllRegionsHandler: no regions to show, err: %s", err)
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			log.Printf("RegionHandlers GetAllRegionsHandler err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong when extracting regions",
			})
			return
		}

		log.Printf("Sent all regions to user with ID: %d", userID)
		c.JSON(http.StatusOK, gin.H{
			"regions": regions,
		})
	}
}
