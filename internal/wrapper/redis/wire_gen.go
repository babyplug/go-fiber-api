// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package redis

import (
	"github.com/google/wire"
	"go-fiber-api/internal/core/config"
)

// Injectors from wire.go:

// Wire is the wire provider for the redis package
func Wire(cfg *config.Configuration) (Client, error) {
	client, err := ProvideClient(cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// wire.go:

var ProviderSet = wire.NewSet(ProvideClient)
