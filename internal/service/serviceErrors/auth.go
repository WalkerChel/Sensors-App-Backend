package serviceErrors

import "errors"

var (
	ErrTokenAlreadyRemoved = errors.New("token has already been deleted from db")
	ErrNoTokenForCheck     = errors.New("token doesn't exist in db")
	ErrParseToken          = errors.New("token parse error")
)
