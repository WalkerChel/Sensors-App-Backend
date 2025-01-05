package serviceErrors

import "errors"

var (
	ErrNoRegionsData = errors.New("there are no records of the regions")
)
