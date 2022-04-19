package launcher

import (
	"math/rand"
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/pkg/config"

	"github.com/A-SoulFan/asasfans-api/internal/pkg/log"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func Run(app fx.Option) {
	fx.New(
		app,

		//全局依赖
		config.Provide(),
		log.Provide(),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func() {
			// default
			rand.Seed(time.Now().UnixNano())
		}),
	).Run()
}
