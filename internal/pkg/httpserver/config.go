package httpserver

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Host               string
	Port               int
	Mode               string
	MaxMultipartMemory uint8
}

func NewConfig(v *viper.Viper, logger *zap.Logger) (*Config, error) {
	var err error
	o := &Config{}
	if err = v.UnmarshalKey("http", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal http option error")
	}

	logger.Info("load http server options success")

	return o, err
}
