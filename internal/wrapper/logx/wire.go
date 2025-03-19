//go:build wireinject
// +build wireinject

//go:generate wire
package logx

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	Provide,
)

func Wire() (*LogX, error) {
	wire.Build(ProviderSet)
	return &LogX{}, nil
}
