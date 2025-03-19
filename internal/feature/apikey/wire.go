//go:build wireinject
// +build wireinject

//go:generate wire
package apikey

import (
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/core/storage/db"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideRepository,

	ProvideService,

	ProvideHandler,
)

func Wire(cfg *config.Configuration, client db.Client) (Handler, error) {
	wire.Build(ProviderSet)

	return &handlerImpl{}, nil
}
