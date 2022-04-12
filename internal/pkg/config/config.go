package config

import (
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(NewViper)
}
