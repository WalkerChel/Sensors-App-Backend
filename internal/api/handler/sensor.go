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
	GetPaginatedSensors(ctx context.Context, limit, offset int64) ([]entities.Sensor, int64, error)
}

type SensorHandlers struct {
	sensorService SensorService
}

func NewSensorHandlers(sensorService SensorService) SensorHandlers {
	return SensorHandlers{
		sensorService: sensorService,
	}
}

// regions/{region_id}/sensors
func (h *SensorHandlers) GetSensorsInRegionHandler(authService ports.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := authService.GetUserIDFromCtx(c, userIDCtxKey)
		if err != nil {
			log.Printf("SensorHandlers GetSensorsInRegionHandler err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong when extracting userID from request's context",
			})
			return
		}

		regIdParam := c.Param(regionIdKey)
		regId, ok := big.NewInt(0).SetString(regIdParam, 10)

		if !ok {
			log.Printf("SensorHandlers GetSensorsInRegionHandler error: %s param is not a number, given param: %s", regionIdKey, regIdParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("given path {%s} param is not a number", regionIdKey),
			})
			return
		}

		if !regId.IsInt64() {
			log.Printf("SensorHandlers GetSensorsInRegionHandler error: %s param is not an int64 type, given param: %s", regionIdKey, regIdParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("given path {%s} param is too large", regionIdKey),
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

// /sensors?page={page}&limit={limit}
func (h *SensorHandlers) GetPaginatedSensorsHandler(authService ports.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := authService.GetUserIDFromCtx(c, userIDCtxKey)
		if err != nil {
			log.Printf("SensorHandlers GetPaginatedSensorsHandler err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong when extracting userID from request's context",
			})
			return
		}

		pageParam := c.Query(pageKey)
		page, ok := big.NewInt(0).SetString(pageParam, 10)

		if !ok || !page.IsUint64() {
			log.Printf("SensorHandlers GetPaginatedSensorsHandler error: wrong {%s} param , given param: %s", pageKey, pageParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("query {%s} param can't be omitted, negative or too large", pageKey),
			})
			return
		}

		pageInt64 := page.Int64()

		limitParam := c.Query(limitKey)
		limit, ok := big.NewInt(0).SetString(limitParam, 10)

		if !ok || !limit.IsInt64() || limit.Int64() < 0 && limit.Int64() != -1 {
			log.Printf("SensorHandlers GetPaginatedSensorsHandler error: wrong {%s} param , given param: %s", limitKey, limitParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("query {%s} param can't be omitted, negative or too large", limitKey),
			})
			return
		}

		limitInt64 := limit.Int64()

		var offset int64

		if limitInt64 == -1 {
			log.Printf("{%s} key is '-1'. Getting all sensors", limitKey)
		} else {
			offset = (pageInt64 - 1) * limitInt64
			log.Printf("calculated offset: (%d[page] - 1) * %d[limit] = %d[offset]", pageInt64, limitInt64, offset)
		}

		sensors, amount, err := h.sensorService.GetPaginatedSensors(c, limitInt64, offset)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrNoSensorsData) {
				log.Printf("SensorHandlers GetPaginatedSensorsHandler: %s", err)
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			log.Printf("SensorHandlers GetPaginatedSensorsHandler error: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong when finding sensors",
			})
			return
		}

		var sensorsResponse apiResponses.SensorsPaginatedResponse

		for _, sensor := range sensors {
			sensorsResponse.Sensors = append(sensorsResponse.Sensors, apiResponses.SensorWithoutRegionResponse{
				ID:        sensor.Id,
				Name:      sensor.Name,
				Longitude: sensor.Longitude,
				Latitude:  sensor.Latitude,
			})
		}

		sensorsResponse.AllSensorsCount = amount

		log.Printf("SensorHandlers GetPaginatedSensorsHandler: successfully sent sensors limit: %d, offset: %d to user with ID: %d", limitInt64, offset, userID)
		c.JSON(http.StatusOK, sensorsResponse)
	}
}
