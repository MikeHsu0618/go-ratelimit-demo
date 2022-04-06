package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	IPLimitPeriod        = 60
	IPLimitMaximum int64 = 60
	script               = `
		local count = redis.call('incr', KEYS[1])
		if count == 1 then
		redis.call('expire', KEYS[1], tonumber(KEYS[2]))
		end
		local reset = redis.call('ttl', KEYS[1])
		return {
			count,
			reset
		}
	`
)

var (
	client *redis.Client
	ctx    = context.Background()
	count  int64
	reset  int64
)

func main() {
	r := gin.Default()
	client = NewRedisClient()
	r.GET("", index)
	err := r.Run()
	if err != nil {
		panic("server error")
	}
}

func index(c *gin.Context) {
	var luaScript = redis.NewScript(script)
	value, err := luaScript.Run(ctx, client, []string{c.ClientIP(), strconv.Itoa(IPLimitPeriod)}).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	result := value.([]interface{})
	count = result[0].(int64)
	reset = result[1].(int64)

	if count > IPLimitMaximum {
		c.JSON(http.StatusTooManyRequests, "too many request, reset: "+strconv.Itoa(int(reset)))
		return
	}

	c.JSON(200, "request:"+strconv.Itoa(int(count))+" reset:"+strconv.Itoa(int(reset)))
}

func NewRedisClient() *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "mypassword",
	})

	pong, err := client.Ping(ctx).Result()
	fmt.Println(pong, err)
	return client
}
