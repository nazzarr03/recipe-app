package config

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
)

func ConnectRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	pong, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		panic("failed to connect to Redis")
	}

	fmt.Println("Redis connected successfully:", pong)
}
