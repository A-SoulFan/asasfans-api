package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

var configFilePath = flag.String("f", "config/config.yml", "set config file which viper will loading.")

func NewViper() *viper.Viper {
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
	return v
}
