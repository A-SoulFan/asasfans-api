package smsclient

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Client interface {
	SendHTML(to, content, subject string) error
}

type SmtpPool struct {
	conf *Config
	pool *email.Pool
}

type Config struct {
	Host               string
	Port               int
	MaxConns           int
	IdleTimeout        int
	PoolWaitTimeout    int
	From               string
	UserName           string
	Password           string
	AuthType           string
	SSL                bool `json:"ssl" yaml:"ssl"`
	InsecureSkipVerify bool
}

func NewConfig(v *viper.Viper, logger *zap.Logger) (*Config, error) {
	var err error
	conf := &Config{}
	if err = v.UnmarshalKey("smsClient", conf); err != nil {
		return nil, errors.Wrap(err, "unmarshal sms-client config error")
	}

	logger.Info("load sms-client options success", zap.String("host", conf.Host), zap.Int("port", conf.Port))

	return conf, err
}

func NewSmtPool(conf *Config) (Client, error) {
	var auth smtp.Auth
	switch conf.AuthType {
	case "login":
		auth = LoginAuth(conf.UserName, conf.Password)
	case "plain":
		auth = smtp.PlainAuth("", conf.UserName, conf.Password, conf.Host)
	default:
		return nil, errors.New("invalid authType, allow: login, plain")
	}

	var pool *email.Pool
	var err error
	var tlsConfig *tls.Config

	if !conf.SSL {
		pool, err = email.NewPool(
			fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			conf.MaxConns,
			auth,
		)
	} else {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: conf.InsecureSkipVerify,
			ServerName:         conf.Host,
		}
		pool, err = email.NewPool(
			fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			conf.MaxConns,
			auth,
			tlsConfig,
		)
	}

	if err != nil {
		return nil, errors.Wrap(err, "creating smtp-pool")
	}

	return &SmtpPool{
		conf: conf,
		pool: pool,
	}, nil
}

func (sp *SmtpPool) SendHTML(to, content, subject string) error {
	e := email.Email{
		From:    sp.conf.From,
		To:      []string{to},
		Subject: subject,
		HTML:    []byte(content),
	}

	err := sp.pool.Send(&e, time.Duration(sp.conf.IdleTimeout)*time.Second)
	if err != nil {
		return errors.Wrap(err, "send email error")
	}

	return nil
}

func Provide() fx.Option {
	return fx.Provide(NewConfig, NewSmtPool)
}
