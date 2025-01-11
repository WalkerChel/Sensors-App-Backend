package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	apiRequests "sensors-app/internal/api/api-requests"
	apiResponses "sensors-app/internal/api/api-responses"
	"sensors-app/internal/api/ports"
	"sensors-app/internal/entities"
	"sensors-app/internal/service/serviceErrors"
	"time"

	"github.com/gin-gonic/gin"
)

type ReadingsService interface {
	FindReadingsBySensorID(ctx context.Context, sensorId int64, time1, time2 time.Time) ([]entities.Reading, error)
}

type ReadingsHandlers struct {
	readingsService ReadingsService
}

func NewReadingsHandlers(readingsService ReadingsService) ReadingsHandlers {
	return ReadingsHandlers{
		readingsService: readingsService,
	}
}

func (h *ReadingsHandlers) GetSensorReadingsHandler(authService ports.Authentication) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := authService.GetUserIDFromCtx(c, userIDCtxKey)
		if err != nil {
			log.Printf("ReadingsHandlers GetSensorReadingsHandler err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "something went wrong when extracting userID from request's context",
			})
			return
		}

		sensIdParam := c.Param(sensorIdKey)
		sensId, ok := big.NewInt(0).SetString(sensIdParam, 10)

		if !ok {
			log.Printf("ReadingsHandlers GetSensorReadingsHandler error: %s param is not a number, given param: %s", regionIdKey, sensIdParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("given path {%s} param is not a number", regionIdKey),
			})
			return
		}

		if !sensId.IsUint64() {
			log.Printf("ReadingsHandlers GetSensorReadingsHandler error: %s param is not an uint64 type, given param: %s", regionIdKey, sensIdParam)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("given path {%s} param is too large or negative", regionIdKey),
			})
			return
		}

		sensIdInt64 := sensId.Int64()

		var dates apiRequests.ReadingsBetweenDates
		var dateStart, dateEnd time.Time

		// can not read body before binding. it closes after first reading
		// reqData, err := c.GetRawData()
		// if err != nil {
		// 	log.Printf("ReadingsHandlers GetSensorReadingsHandler get request data err: %s", err)
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		// 		"error": "something went wrong when reading request's data",
		// 	})
		// 	return
		// }

		if c.Request.Body != nil { // && string(reqData) != "" {
			if err := c.ShouldBindJSON(&dates); err != nil {
				if !errors.Is(err, io.EOF) {
					log.Printf("ReadingsHandlers GetSensorReadingsHandler bindJSON err: %s", err)
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
						"error": "check request's body",
					})
					return
				}
				log.Print("ReadingsHandlers GetSensorReadingsHandler bindJSON EOF, trying to get all readings")

			}
			if dates.StartDate != nil {
				dateStart = time.Date(dates.StartDate.Year,
					time.Month(dates.StartDate.Month),
					dates.StartDate.Day,
					dates.StartDate.Hour,
					dates.StartDate.Minute,
					dates.StartDate.Second, 0, time.UTC)
			}

			if dates.EndDate != nil {
				dateEnd = time.Date(dates.EndDate.Year,
					time.Month(dates.EndDate.Month),
					dates.EndDate.Day,
					dates.EndDate.Hour,
					dates.EndDate.Minute,
					dates.EndDate.Second, 0, time.UTC)
			} else if dates.StartDate != nil && dates.EndDate == nil {
				dateEnd = time.Now()
			}
		}

		readings, err := h.readingsService.FindReadingsBySensorID(c, sensIdInt64, dateStart, dateEnd)
		if err != nil {
			if errors.Is(err, serviceErrors.ErrNoReadingsData) {
				log.Printf("ReadingsHandlers GetSensorReadingsHandler no data: %s", err)
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
			if errors.Is(err, serviceErrors.ErrIncorrectDates) {
				log.Printf("ReadingsHandlers GetSensorReadingsHandler dates err: %s", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "start date can not be after end date",
				})
				return
			}
			if errors.Is(err, serviceErrors.ErrEndDateAfterCurrDate) {
				log.Printf("ReadingsHandlers GetSensorReadingsHandler dates err: %s", err)
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "end date can not be after current date",
				})
				return
			}
			log.Printf("ReadingsHandlers GetSensorReadingsHandler getting readings err: %s", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong when finding readings",
			})
			return
		}

		var readingsPlot []apiResponses.ReadingForPlotResponse

		for _, reading := range readings {
			readungPlot := apiResponses.ReadingForPlotResponse{
				Temperature: reading.Temperature,
				CreatedAt:   reading.CreatedAt,
			}
			readingsPlot = append(readingsPlot, readungPlot)
		}

		log.Printf("ReadingsHandlers GetSensorReadingsHandler: successfully sent readings from sensor: %d to userID: %d", sensIdInt64, userID)
		c.JSON(http.StatusOK, apiResponses.ReadingsForSensorResponse{
			Readings: readingsPlot,
			SensorId: sensIdInt64,
		})
	}
}
