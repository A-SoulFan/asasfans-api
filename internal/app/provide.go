package app

import (
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/middlewares"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/router"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/service"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/util/password"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/cache"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/smsclient"

	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Options(
		router.Provide(),
		MiddlewareProvider(),
		ServiceProvider(),

		smsclient.Provide(),
		fx.Provide(cache.NewGoCache),
		fx.Provide(smsclient.NewRandomNumberCodeGenerator),
		fx.Provide(password.NewDefaultPasswordHandler),
	)
}

func MiddlewareProvider() fx.Option {
	return fx.Provide(
		middlewares.NewErrorInterceptor,
		middlewares.NewSession,
	)
}

func ServiceProvider() fx.Option {
	return fx.Provide(
		service.NewBilbilVideo,
		service.NewAuth,
		service.NewUser,
	)
}
