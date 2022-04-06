package app

import (
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/router"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/service"

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
