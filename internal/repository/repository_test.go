package repository

import (
	"log"
	"os"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"ratelimit/internal/lua"
)

var (
	client *redis.Client
	repo   Repository
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	repo = NewRepository(client)
	code := m.Run()
	os.Exit(code)
}

func TestGetCountAndTTL(t *testing.T) {
	IP := "127.0.0.1"
	for i := 1; i <= 10; i++ {
		count, ttl, err := repo.GetCountAndTTL(IP)

		if count == 1 {
			assert.Equal(t, ttl, int64(lua.IPLimitPeriod))
		}
		assert.Equal(t, count, int64(i))
		assert.Equal(t, err, nil)
	}
}
