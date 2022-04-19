package repository

import (
	"fmt"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	userTableName = "user"
)

func NewUser(tx *gorm.DB) idl.UserRepository {
	return NewUserMysqlImpl(tx)
}

type UserMysqlImpl struct {
	tx *gorm.DB
}

func NewUserMysqlImpl(tx *gorm.DB) idl.UserRepository {
	return &UserMysqlImpl{tx: tx}
}

func (impl *UserMysqlImpl) Create(u *idl.User) (isNew bool, err error) {
	result := impl.tx.Table(userTableName).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}, {Name: "email"}},
		UpdateAll: true,
	}).Create(&u)

	if result.Error != nil {
		return false, errors.Wrap(result.Error, fmt.Sprintf("insert %s error", userTableName))
	}

	if u.Id == 0 {
		impl.tx.Table(userTableName).
			Where("email = ?", u.Email).
			Find(&u)
		isNew = false
	} else {
		isNew = true
	}

	return isNew, nil
}

func (impl *UserMysqlImpl) Update(u *idl.User, fields ...string) error {
	var tx *gorm.DB
	if len(fields) > 0 {
		tx = impl.tx.Table(userTableName).Select(fields)
	}

	result := tx.Updates(&u)
	if result.Error != nil {
		return errors.Wrap(result.Error, fmt.Sprintf("update %s error", userTableName))
	}

	return nil
}

func (impl *UserMysqlImpl) FindUserByID(uid uint64) (*idl.User, error) {
	var u *idl.User
	result := impl.tx.Table(userTableName).
		Where("id = ?", uid).
		Find(&u)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("select %s error", userTableName))
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return u, nil
}

func (impl *UserMysqlImpl) FindUserByEmail(email string) (*idl.User, error) {
	var u *idl.User
	result := impl.tx.Table(userTableName).
		Where("email = ?", email).
		Find(&u)

	if result.Error != nil {
		return nil, errors.Wrap(result.Error, fmt.Sprintf("select %s error", userTableName))
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	return u, nil
}
