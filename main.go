package main

import (
	"github.com/gin-gonic/gin"
	"ratelimit/internal/repository"
	"ratelimit/internal/service"
	"ratelimit/pkg/redislib"
)

func main() {
	r := SetUpRouter()
	err := r.Run()
	if err != nil {
		panic("server error")
	}
}

func SetUpRouter() *gin.Engine {
	r := gin.Default()
	client := redislib.NewRedisClient()
	repo := repository.NewRepository(client)
	svc := service.NewService(repo)
	r.GET("", svc.RateLimit)
	return r
}
