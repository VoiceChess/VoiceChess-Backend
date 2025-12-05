package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func New() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "redis://default:EltyWJRPuTLqwkuTLOuAehtHauvFlmGB@redis.railway.internal:6379",
		Password: "",
		DB:       0,
	})
}

func RedisSet(key, value string) error {
	return Rdb.Set(Ctx, key, value, 0).Err()
}

func RedisGet(key string) (string, error) {
	return Rdb.Get(Ctx, key).Result()
}
