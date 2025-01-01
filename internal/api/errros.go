package api

import "errors"

var (
	ErrInsufficientFields = errors.New("Not enough fields in request body")
)
