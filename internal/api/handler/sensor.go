package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	apiResponses "sensors-app/internal/api/api-responses"
	"sensors-app/internal/api/ports"
	"sensors-app/internal/entities"
	"sensors-app/internal/service/serviceErrors"

	"github.com/gin-gonic/gin"
)

type SensorService interface {
	GetSensorsByRegionID(ctx context.Context, regionId int64) ([]entities.Sensor, error)
}

type SensorHandlers struct {
	sensorService SensorService
}

func NewSensorHandlers(sensorService SensorService) SensorHandlers {
	return SensorHandlers{
		sensorService: sensorService,
	}
}

// sensors/query?regionID={regionID}
func (h *SensorHandlers) GetSensorsInRegionHandler(authService ports.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := authService.GetUserIDFromCtx(c, userIDCtxKey)
		if err != nil {
			log.Printf("RegionHandlers GetAllRegionsHandler err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong when extracting userID from request's context",
			})
			return
		}

		regIdParam := c.Query(regionIdQuery)
		regId, ok := big.NewInt(0).SetString(regIdParam, 10)

		if !ok {
			log.Printf("SensorHandlers GetSensorsInRegionHandler error: %s param is not a number, given param: %s", regionIdQuery, regIdParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("given %s param is not a number", regionIdQuery),
			})
			return
		}

		if !regId.IsInt64() {
			log.Printf("SensorHandlers GetSensorsInRegionHandler error: %s param is not an int64 type, given param: %s", regionIdQuery, regIdParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("given %s param is too large", regionIdQuery),
			})
			return
		}

		regIdInt64 := regId.Int64()

		sensors, err := h.sensorService.GetSensorsByRegionID(c, regIdInt64)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrNoSensorsData) {
				log.Printf("SensorHandlers GetSensorsInRegionHandler: %s", err)
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			log.Printf("SensorHandlers GetSensorsInRegionHandler error: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong when finding sensors",
			})
			return
		}

		var sensorsResponse apiResponses.SensorsInRegionResponse
		var sensorsInRegion []apiResponses.SensorWithoutRegionResponse

		for _, sensor := range sensors {
			sensResp := apiResponses.SensorWithoutRegionResponse{
				ID:        sensor.Id,
				Name:      sensor.Name,
				Longitude: sensor.Longitude,
				Latitude:  sensor.Latitude,
			}
			sensorsInRegion = append(sensorsInRegion, sensResp)
		}

		sensorsResponse.Sensors = sensorsInRegion
		sensorsResponse.RegionID = regIdInt64

		log.Printf("SensorHandlers GetSensorsInRegionHandler: successfully sent sensors from region #%d to user with ID: %d", regIdInt64, userID)
		c.JSON(http.StatusOK, sensorsResponse)

	}
}
