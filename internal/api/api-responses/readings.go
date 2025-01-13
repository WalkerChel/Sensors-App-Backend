package apiResponses

import "time"

type ReadingForPlotResponse struct {
	Temperature float64   `json:"reading_temperature"`
	CreatedAt   time.Time `json:"reading_created_at"`
}

type ReadingsForSensorResponse struct {
	Readings       []ReadingForPlotResponse `json:"readings"`
	SensorId       int64                    `json:"sensor_id"`
	MinTemperature float64                  `json:"temperature_min"`
	AvgTemperature float64                  `json:"temperature_avg"`
	MaxTemperature float64                  `json:"temperature_max"`
}
