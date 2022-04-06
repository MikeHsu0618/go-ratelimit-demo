package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ratelimit/internal/lua"
	"ratelimit/internal/repository"
)

type Service interface {
	RateLimit(c *gin.Context)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s service) RateLimit(c *gin.Context) {
	count, ttl, err := s.repo.GetCountAndTTL(c.ClientIP())
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	if count > lua.IPLimitMaximum {
		c.JSON(http.StatusTooManyRequests, "too many request, ttl: "+strconv.Itoa(int(ttl)))
		return
	}

	c.JSON(http.StatusOK, "request:"+strconv.Itoa(int(count))+" ttl:"+strconv.Itoa(int(ttl)))
}
