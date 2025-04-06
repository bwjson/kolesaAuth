package redis

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

func NewRedisClient(ctx context.Context, addr, user, pass string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: pass,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Ошибка при подключении к Redis:", err)
	}

	return rdb
}
