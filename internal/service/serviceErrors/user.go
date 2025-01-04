package serviceErrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user alredy exists")
	ErrNoUserInfo        = errors.New("user doesn't exist")
)

var (
	ErrNoUserIDInCtx      = errors.New("no user id found in request's context")
	ErrUserIDNotInt64Type = errors.New("user id value in request's context is not int64 type")
)
