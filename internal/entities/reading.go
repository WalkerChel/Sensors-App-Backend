package entities

import "time"

type Reading struct {
	Id          int64     `json:"reading_id" db:"id"`
	SensorId    int64     `json:"sensor_id" db:"sensor_id"`
	Temperature float64   `json:"reading_temperature" db:"temperature"`
	CreatedAt   time.Time `json:"reading_created_at" db:"created_at"`
}
