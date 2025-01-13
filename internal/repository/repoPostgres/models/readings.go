package models

import "sensors-app/internal/entities"

type ReadingsWthStats struct {
	entities.Reading
	MinTemperature float64 `db:"temperature_min"`
	AvgTemperature float64 `db:"temperature_avg"`
	MaxTemperature float64 `db:"temperature_max"`
}
