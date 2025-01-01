package repoErrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user alredy exists")
)
