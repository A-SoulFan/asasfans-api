package main

import (
	"context"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/spider"
	"github.com/A-SoulFan/asasfans-api/internal/launcher"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilbil"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/database"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func main() {
	launcher.Run(func(viper *viper.Viper) fx.Option {
		return fx.Options(
			database.Provide(),
			fx.Provide(spider.NewVideo),
			fx.Provide(bilbil.NewSDK),
			fx.Invoke(func(lifecycle fx.Lifecycle, spiderVideo *spider.Video) {
				lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						spiderVideo.Run()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return spiderVideo.Stop()
					},
				})
			}),
		)
	})
}