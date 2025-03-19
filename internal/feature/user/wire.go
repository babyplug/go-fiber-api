//go:build wireinject
// +build wireinject

//go:generate wire
package user

import (
	"go-fiber-api/internal/core/storage/db"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideRepository,

	ProvideService,

	ProvideHandler,
)

func Wire(db db.Client) (Handler, error) {
	wire.Build(ProviderSet)

	return &handlerImpl{}, nil
}
