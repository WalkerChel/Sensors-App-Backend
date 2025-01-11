package apiRequests

import (
	"sensors-app/internal/entities"
)

type ReadingsBetweenDates struct {
	StartDate *entities.DateTime `json:"start_date" binding:"required"`
	EndDate   *entities.DateTime `json:"end_date" binding:"required"`
}
