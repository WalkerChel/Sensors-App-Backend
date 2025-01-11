package serviceErrors

import "errors"

var (
	ErrNoReadingsData       = errors.New("there are no records of the readings")
	ErrIncorrectDates       = errors.New("start date can't be after end date")
	ErrEndDateAfterCurrDate = errors.New("end date can't be after current date")
)
