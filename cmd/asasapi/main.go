package main

import (
	"context"

	"github.com/A-SoulFan/asasfans-api/internal/app"
	"github.com/A-SoulFan/asasfans-api/internal/launcher"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/database"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/httpserver"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func main() {
	launcher.Run(func(viper *viper.Viper) fx.Option {
		return fx.Options(
			database.Provide(),
			httpserver.Provide(),
			app.Provide(),
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
