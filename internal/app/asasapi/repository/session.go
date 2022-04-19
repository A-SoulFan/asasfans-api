package repository

import (
	"fmt"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	sessionTableName = "session"
)

func NewTokenStorage(tx *gorm.DB) idl.TokenStorage {
	return NewSessionMysqlImpl(tx)
}

type SessionMysqlImpl struct {
	tx *gorm.DB
}

func NewSessionMysqlImpl(tx *gorm.DB) idl.TokenStorage {
	return &SessionMysqlImpl{tx: tx}
}

func (impl *SessionMysqlImpl) Get(key string) (token *idl.Token, err error) {
	token = &idl.Token{}
	result := impl.tx.Table(sessionTableName).Where("session_key = ?", key).Select("u_id").Find(&token.UserId)
	if result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("select %s fail", sessionTableName))
	}

	return token, nil
}

func (impl *SessionMysqlImpl) Set(key string, token *idl.Token) error {
	result := impl.tx.Table(sessionTableName).Clauses(clause.OnConflict{DoNothing: true}).Create(struct {
		UserId     uint64 `gorm:"column:u_id"`
		SessionKey string
	}{
		UserId:     token.UserId,
		SessionKey: key,
	})

	if result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf("insert %s fail", sessionTableName))
	}

	if result.RowsAffected == 0 {
		return &idl.UniqueConflictError{Message: "session_key unique conflict"}
	}

	return nil
}

func (impl *SessionMysqlImpl) Del(key string) error {
	result := impl.tx.Table(sessionTableName).Where("session_key = ?", key).Delete(nil)
	if result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf("delete %s fail", sessionTableName))
	}
	return nil
}

func (impl *SessionMysqlImpl) ClearTokens(token *idl.Token) error {
	result := impl.tx.Table(sessionTableName).Where("u_id = ?", token.UserId).Delete(nil)
	if result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf("delete %s fail", sessionTableName))
	}
	return nil
}
