package ports

import (
	"context"
	"sensors-app/internal/entities"
)

type Authentication interface {
	CheckToken(cxt context.Context, userId int64) (bool, error)
	CreateToken(cxt context.Context, userId int64, cnf entities.JWT) (string, error)
	DeleteToken(cxt context.Context, userId int64, token string) error
}
