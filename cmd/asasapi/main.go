package main

import (
	"asasfans/internal/app/asasapi/router"
	"asasfans/internal/launcher"
	"asasfans/internal/pkg/database"
	"asasfans/internal/pkg/httpserver"
	"context"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func main() {
	launcher.Run(func(viper *viper.Viper) fx.Option {
		return fx.Options(
			database.Provide(),
			httpserver.Provide(),
			fx.Provide(router.InitRouters),
			fx.Invoke(func(lifecycle fx.Lifecycle, ginServer *httpserver.Server) {
				lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return ginServer.Start()
					},
					OnStop: func(ctx context.Context) error {
						return ginServer.Stop()
					},
				})
			}),
		)
	})
}
