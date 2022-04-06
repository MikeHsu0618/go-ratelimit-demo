package repository

import (
	"context"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"ratelimit/internal/lua"
)

type Repository interface {
	GetCountAndTTL(key string) (int64, int64, error)
}

type repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) Repository {
	return &repository{client}
}

func (r repository) GetCountAndTTL(key string) (count int64, ttl int64, err error) {
	var luaScript = redis.NewScript(lua.Script)
	value, err := luaScript.Run(context.Background(), r.client, []string{key, strconv.Itoa(lua.IPLimitPeriod)}).Result()
	if err != nil {
		log.Fatalln("redis error")
		return 0, 0, err
	}

	result := value.([]interface{})
	count = result[0].(int64)
	ttl = result[1].(int64)
	return
}
