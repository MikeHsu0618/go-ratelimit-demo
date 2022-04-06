package service

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"ratelimit/internal/repository"
)

var (
	client *redis.Client
	svc    Service
	r      *gin.Engine
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	r = gin.Default()
	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	repo := repository.NewRepository(client)
	svc = NewService(repo)
	r.GET("/", svc.RateLimit)

	code := m.Run()
	os.Exit(code)
}

func TestRateLimit(t *testing.T) {
	for i := 1; i <= 100; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if i <= 60 {
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}
}
