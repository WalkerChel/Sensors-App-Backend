package db

import (
	"context"
	"net"

	"github.com/redis/go-redis/v9"
)

func NewRedisDB(ctx context.Context, host, port, password string, DB int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(host, port),
		Password: password,
		DB:       DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
