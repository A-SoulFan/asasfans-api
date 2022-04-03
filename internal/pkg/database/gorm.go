package database

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Options is  configuration of database
type Options struct {
	DSN                string `yaml:"dsn"`
	Debug              bool
	SetMaxIdleConns    int
	SetMaxOpenConns    int
	SetConnMaxLifetime int
}

func NewOptions(v *viper.Viper, logger *zap.Logger) (*Options, error) {
	var err error
	o := &Options{}
	if err = v.UnmarshalKey("db", o); err != nil {
		return nil, errors.Wrap(err, "unmarshal db option error")
	}

	logger.Info("load database options success", zap.String("dsn", o.DSN))

	return o, err
}

// NewDatabase 初始化数据库
func NewDatabase(o *Options) (db *gorm.DB, err error) {
	var dialector gorm.Dialector
	dialector = mysql.Open(o.DSN)

	db, err = gorm.Open(dialector, &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 nil,
	})
	if err != nil {
		return nil, err
	}

	if sqlDb, err := db.DB(); err != nil {
		return nil, err
	} else {
		sqlDb.SetMaxIdleConns(o.SetMaxIdleConns)
		sqlDb.SetMaxOpenConns(o.SetMaxOpenConns)
		sqlDb.SetConnMaxLifetime(time.Duration(o.SetConnMaxLifetime) * time.Minute)
	}

	if o.Debug {
		db = db.Debug()
	}

	return db, nil
}

func Provide() fx.Option {
	return fx.Provide(NewOptions, NewDatabase)
}
