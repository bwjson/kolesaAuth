package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
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
	return nil
}

func (r *Repository) Delete(ctx context.Context, phoneNumber string) error {
	return nil
}

func (r *Repository) IsValidCode(ctx context.Context, phoneNumber, verificationCode string) (bool, error) {
	return false, nil
}
