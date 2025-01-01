package repoRedis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepo struct {
	db *redis.Client
}

func NewTokenRepo(db *redis.Client) TokenRepo {
	return TokenRepo{
		db: db,
	}
}

func (r *TokenRepo) StoreToken(cxt context.Context, userId int64, token string, tokenTTL time.Duration) error {
	key := fmt.Sprintf("userId:%d:token", userId)

	if err := r.db.Set(cxt, key, token, tokenTTL).Err(); err != nil {
		log.Printf("RedisRepo StoreToken err: %s", err)
		return err
	}

	return nil
}

func (r *TokenRepo) DeleteToken(cxt context.Context, userId int64) error {
	key := fmt.Sprintf("userId:%d:token", userId)

	if err := r.db.Del(cxt, key).Err(); err != nil {
		log.Printf("RedisRepo DeleteToken err: %s", err)
		return err
	}

	return nil
}

func (r *TokenRepo) TokenExists(cxt context.Context, userId int64) (bool, error) {
	key := fmt.Sprintf("userId:%d:token", userId)

	_, err := r.db.Get(cxt, key).Result()

	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get value, error: %w", err)
	}

	return true, nil
}
