package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

//type Repository interface {
//	Save(ctx context.Context, phoneNumber, verificationCode string) error
//	Delete(ctx context.Context, phoneNumber string) error
//	IsValidCode(ctx context.Context, phoneNumber, verificationCode string) (bool, error)
//}

type Repository struct {
	cache *redis.Client
}

func NewRepository(cache *redis.Client) *Repository {
	return &Repository{cache}
}

func (r *Repository) Save(ctx context.Context, key, value string, ttl time.Duration) error {
	// TODO: Change ttl param using env variable
	err := r.cache.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, key string) error {
	err := r.cache.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return val, nil
}

func (r *Repository) IsValidCode(ctx context.Context, phoneNumber, verificationCode string) (bool, error) {
	val, err := r.cache.Get(ctx, phoneNumber).Result()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if val != verificationCode {
		return false, nil
	}

	return true, nil
}
