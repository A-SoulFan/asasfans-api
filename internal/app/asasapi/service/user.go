package service

import (
	"context"

	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser() *User {
	return &User{}
}

func (u *User) Info(ctx context.Context) {

}
