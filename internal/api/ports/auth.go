package ports

import (
	"context"
	"sensors-app/internal/entities"
)

type Authentication interface {
	CheckToken(cxt context.Context, userId int64, token string) (bool, error)
	ParseToken(token string, cnf entities.JWT) (int64, error)
	CreateToken(cxt context.Context, userId int64, cnf entities.JWT) (string, error)
	DeleteToken(cxt context.Context, userId int64) error
	GetUserIDFromCtx(ctx context.Context, key string) (int64, error)
}
