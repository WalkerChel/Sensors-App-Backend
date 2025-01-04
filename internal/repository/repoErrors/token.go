package repoErrors

import "errors"

var (
	ErrNoToken = errors.New("token doesn't exist in db")
)
