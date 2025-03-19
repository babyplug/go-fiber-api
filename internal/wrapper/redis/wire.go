//go:build wireinject
// +build wireinject

//go:generate wire
package redis

import (
	"go-fiber-api/internal/core/config"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(ProvideClient)

// Wire is the wire provider for the redis package
func Wire(cfg *config.Configuration) (Client, error) {
	wire.Build(ProviderSet)

	return &clientImpl{}, nil
}
