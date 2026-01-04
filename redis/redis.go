package redis

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func New() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
}

func RedisSet(key, value string) error {
	return Rdb.Set(Ctx, key, value, 0).Err()
}

func RedisGet(key string) (string, error) {
	return Rdb.Get(Ctx, key).Result()
}
