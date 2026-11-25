package apiRequests

import (
	"github.com/WalkerChel/sensors-api/internal/entities"
)

type ReadingsBetweenDates struct {
	StartDate *entities.DateTime `json:"start_date" binding:"required"`
	EndDate   *entities.DateTime `json:"end_date" binding:"required"`
}
