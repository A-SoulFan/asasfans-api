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
	launcher.Run(newMigrator())
}

func newMigrator() fx.Option {
	return fx.Options(
		database.Provide(),
		fx.Provide(spider.NewDbMigrate),
		fx.Provide(bilibili.NewSDK),
		fx.Invoke(lc),
	)
}

func lc(lifecycle fx.Lifecycle, dbmigrate *spider.DBMigrate) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return dbmigrate.Run(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return dbmigrate.Stop(ctx)
		},
	})
}
