//go:build wireinject
// +build wireinject

//go:generate wire
package apikey

import (
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/feature/apikey"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Provide,
)

func Wire(config *config.Configuration, apikeyService apikey.Service) (Middleware, error) {
	wire.Build(ProviderSet)

	return &middlewareImpl{}, nil
}
