package repoErrors

import "errors"

var (
	ErrNoRecords = errors.New("no data in db")
)
