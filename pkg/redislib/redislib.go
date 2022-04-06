package redislib

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "mypassword",
	})

	pong, err := client.Ping(context.Background()).Result()
	fmt.Println(pong, err)
	return client
}
