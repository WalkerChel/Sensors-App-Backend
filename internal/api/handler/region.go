package handler

import (
	"context"
	"errors"
	"net/http"
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

func (h *RegionHandlers) GetAllRegions() gin.HandlerFunc {
	return func(c *gin.Context) {
		regions, err := h.regionService.GetAllRegions(c)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrNoRegionsData) {
				c.JSON(http.StatusNoContent, gin.H{
					"message": "No regions records are found",
				})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong when extracting regions",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"regions": regions,
		})
	}
}
