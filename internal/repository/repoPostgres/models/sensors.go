package models

import "sensors-app/internal/entities"

type SensorsWithAllAmount struct {
	entities.Sensor
	Amount int64 `db:"all_sensors_amount"`
}
