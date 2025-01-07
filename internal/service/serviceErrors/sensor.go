package serviceErrors

import "errors"

var (
	ErrNoSensorsData = errors.New("there are no records of the sensors")
)
