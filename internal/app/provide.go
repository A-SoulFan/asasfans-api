package app

import (
	"asasfans/internal/app/asasapi/router"
	"asasfans/internal/app/asasapi/service"

	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Options(
		router.Provide(),
		ServiceProvider(),
	)
}

func ServiceProvider() fx.Option {
	return fx.Provide(service.NewBilbilVideo)
}
