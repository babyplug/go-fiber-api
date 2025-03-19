//go:build wireinject
// +build wireinject

//go:generate wire
package config

import (
	"go-fiber-api/internal/wrapper/logx"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Provide,
)

func Wire(logx *logx.LogX) (*Configuration, error) {
	wire.Build(ProviderSet)
	return &Configuration{}, nil
}
