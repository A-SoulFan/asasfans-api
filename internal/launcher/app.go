package launcher

import (
	"asasfans/internal/pkg/log"
	"flag"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var configFilePath = flag.String("f", "config/config.yml", "set config file which viper will loading.")

func Run(build func(*viper.Viper) fx.Option) {
	flag.Parse()

	var (
		err error
		v   = viper.New()
	)
	v.AddConfigPath(".")
	v.SetConfigFile(*configFilePath)

	if err = v.ReadInConfig(); err == nil {
		fmt.Printf("use config file -> %s\n", v.ConfigFileUsed())
	} else {
		panic(err)
	}

	fx.New(
		build(v),
		fx.Provide(func() *viper.Viper {
			return v
		}),
		log.Provide(),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
	).Run()
}
