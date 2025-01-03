package serviceErrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user alredy exists")
	ErrNoUserInfo        = errors.New("user doesn't exist")
)
