package repoRedis

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sensors-app/internal/repository/repoErrors"
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

	log.Printf("Stored key: '%s' in redis db", key)

	return nil
}

func (r *TokenRepo) DeleteToken(cxt context.Context, userId int64) error {
	key := fmt.Sprintf("userId:%d:token", userId)

	if err := r.db.Del(cxt, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("RedisRepo DeleteToken err: %s", err)
			return repoErrors.ErrNoToken
		}
		log.Printf("RedisRepo DeleteToken err: %s", err)
		return err
	}

	return nil
}

func (r *TokenRepo) GetTokenByUserID(cxt context.Context, userId int64) (string, error) {
	key := fmt.Sprintf("userId:%d:token", userId)

	token, err := r.db.Get(cxt, key).Result()

	log.Printf("REDIS Getting key: %s", key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("RedisRepo TokenExists err: %s", err)
			return "", repoErrors.ErrNoToken
		}
		log.Printf("RedisRepo TokenExists err: %s", err)
		return "", fmt.Errorf("failed to get value, error: %w", err)

	}

	return token, nil
}
