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

func (r *Repository) Save(ctx context.Context, phoneNumber, verificationCode string) error {
	// TODO: Change ttl param using env variable
	err := r.cache.Set(ctx, phoneNumber, verificationCode, 3*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, phoneNumber string) error {
	err := r.cache.Del(ctx, phoneNumber).Err()
	if err != nil {
		return err
	}
	return nil
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
