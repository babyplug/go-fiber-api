//go:build wireinject
// +build wireinject

//go:generate wire
package app

import (
	"go-fiber-api/internal/core/config"
	apikey_middleware "go-fiber-api/internal/core/middleware/apikey"
	"go-fiber-api/internal/core/middleware/cache"
	"go-fiber-api/internal/core/storage/db"
	"go-fiber-api/internal/wrapper/logx"
	"go-fiber-api/internal/wrapper/redis"

	"go-fiber-api/internal/feature/apikey"
	"go-fiber-api/internal/feature/user"

	"github.com/go-resty/resty/v2"
	"github.com/google/wire"
)

func New(client *resty.Client) (*Application, error) {
	wire.Build(
		Provide,
		config.ProviderSet,
		logx.ProviderSet,
		db.ProviderSet,
		user.ProviderSet,
		redis.ProviderSet,
		cache.ProviderSet,
		apikey.ProviderSet,
		apikey_middleware.ProviderSet,
	)

	return &Application{}, nil
}
