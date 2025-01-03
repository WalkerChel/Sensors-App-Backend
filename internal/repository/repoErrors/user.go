package repoErrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user alredy exists")
	ErrNoUser            = errors.New("User doesn't exist")
)
