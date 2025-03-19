package cache

import (
	"fmt"
	"go-fiber-api/internal/wrapper/redis"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

const (
	baseKey string = "cache:%s:%s:%s"
)

var (
	m     *cacheMiddlewareImpl
	mOnce sync.Once
)

type CacheMiddleware interface {
	RedisCacheMiddleware(ttl int) fiber.Handler
}

type cacheMiddlewareImpl struct {
	rc redis.Client
}

func New(rc redis.Client) CacheMiddleware {
	mOnce.Do(func() {
		m = &cacheMiddlewareImpl{
			rc: rc,
		}
	})

	return m
}

func Close() {
	mOnce = sync.Once{}
}

func (m *cacheMiddlewareImpl) RedisCacheMiddleware(ttl int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Clear cache for POST requests
		if c.Method() == fiber.MethodPost {
			if strings.Contains(c.Path(), "/service") {
				return c.Next()
			}
			cacheKey := fmt.Sprintf(baseKey, fiber.MethodGet, c.Path(), "")
			err := m.rc.Del(cacheKey)
			if err != nil {
				logrus.Warnf("failed to clear cache: %v", err)
			}
			return c.Next()
		}

		// Only cache GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		// Generate cache key using method + path
		queries := ""
		if len(c.Queries()) > 0 {
			queries = string(c.Context().URI().QueryString())
		}
		cacheKey := fmt.Sprintf(baseKey, c.Method(), c.Path(), queries)

		// var cachedResponse model.ResponseDTO

		// err := m.rc.Get(cacheKey, &cachedResponse)

		// if err == nil {
		// 	// Return cached response
		// 	c.Set("X-Cache", "HIT") // Mark response as cache hit

		// 	return c.JSON(cachedResponse)
		// } else {
		// 	// Mark as cache miss
		// 	c.Set("X-Cache", "MISS")
		// }

		// Continue with request processing
		err := c.Next()
		if err != nil {
			return err
		}

		// Store response in Redis
		responseBody := c.Response().Body()
		if len(responseBody) > 0 {
			err := m.rc.Set(cacheKey, responseBody, ttl)
			if err != nil {
				logrus.Warnf("failed to set cache: %v", err)
				return nil
			}

			logrus.Warnf("success to set cache: %v", err)
		}

		return nil
	}
}
