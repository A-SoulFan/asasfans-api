package main

import (
	"context"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/spider"
	"github.com/A-SoulFan/asasfans-api/internal/launcher"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/bilibili"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/database"
	"go.uber.org/fx"
)

func main() {
	launcher.Run(newSpider())
}

func newSpider() fx.Option {
	return fx.Options(
		database.Provide(),
		fx.Provide(spider.NewVideo),
		fx.Provide(spider.NewUpdate),
		fx.Provide(bilibili.NewSDK),
		fx.Invoke(lc),
	)
}

func lc(lifecycle fx.Lifecycle, spiderVideo *spider.Video, spiderUpdate *spider.Update, shutdown fx.Shutdowner) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return spiderVideo.Run(ctx)
		},
		OnStop: func(ctx context.Context) error {
			if err := spiderVideo.Stop(ctx); err != nil {
				return err
			}
			return shutdown.Shutdown()
		},
	})

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return spiderUpdate.Run(ctx)
		},
		OnStop: func(ctx context.Context) error {
			if err := spiderUpdate.Stop(ctx); err != nil {
				return err
			}
			return shutdown.Shutdown()
		},
	})
}
