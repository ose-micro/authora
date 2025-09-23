package cache

import (
	"github.com/ose-micro/authora/internal/business/user"
	"github.com/ose-micro/authora/internal/infrastruture/cache/token"
	"github.com/ose-micro/core/logger"
	"github.com/ose-micro/core/tracing"
	"github.com/ose-micro/redis"
)

type Cache struct {
	Token user.Cache
}

func Inject(client *redis.Client, log logger.Logger, trace tracing.Tracer) *Cache {
	return &Cache{
		Token: token.NewTokenCache(client, log, trace),
	}
}
