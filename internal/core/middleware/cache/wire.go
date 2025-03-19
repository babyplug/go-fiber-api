//go:build wireinject
// +build wireinject

//go:generate wire
package cache

import (
	"go-fiber-api/internal/wrapper/redis"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	New,
)

func Wire(client redis.Client) CacheMiddleware {
	wire.Build(ProviderSet)

	return &cacheMiddlewareImpl{}
}
