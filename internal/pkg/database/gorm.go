package database

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Options is  configuration of database
type Options struct {
	Type               string
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
	dialector, err := newDialector(o.Type, o.DSN)

	if err != nil {
		return nil, err
	}

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

func newDialector(t, dsn string) (gorm.Dialector, error) {
	switch strings.ToLower(t) {
	case "sqlite":
		return sqlite.Open(dsn), nil
	case "mysql":
		return mysql.Open(dsn), nil
	default:
		return nil, errors.New("unsupported database type")
	}
}

//DB module 提供DB 和 DBConfig

func Provide() fx.Option {
	return fx.Provide(NewOptions, NewDatabase)
}
