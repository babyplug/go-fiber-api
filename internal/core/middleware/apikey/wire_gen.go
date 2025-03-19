// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package apikey

import (
	"github.com/google/wire"
	"go-fiber-api/internal/core/config"
	"go-fiber-api/internal/feature/apikey"
)

// Injectors from wire.go:

func Wire(config2 *config.Configuration, apikeyService apikey.Service) (Middleware, error) {
	middleware := Provide(config2, apikeyService)
	return middleware, nil
}

// wire.go:

var ProviderSet = wire.NewSet(
	Provide,
)
