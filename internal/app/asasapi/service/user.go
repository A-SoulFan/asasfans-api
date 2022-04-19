package service

import (
	"context"
	"strings"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/help"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	allowUpdateField = "avatar,cover,nickname,gender,short_desc"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) Info(ctx context.Context) (*idl.UserInfoResp, error) {
	ut := help.GetUserToken(ctx)
	if ut == nil {
		return nil, apperrors.NewServiceError(NotFoundCtxUserToken, "get user token fail")
	}

	tx := u.db.WithContext(ctx)
	userRepo := repository.NewUser(tx)
	user, err := userRepo.FindUserByID(ut.UserId)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, apperrors.NewServiceError(InvalidUserId, "get user by id error")
	}

	return idl.NewUserInfoResp(user), nil
}

func (u *User) Update(ctx context.Context, req idl.UserUpdateReq) error {
	ut := help.GetUserToken(ctx)
	if ut == nil {
		return apperrors.NewServiceError(NotFoundCtxUserToken, "get user token fail")
	}

	tx := u.db.WithContext(ctx)
	userRepo := repository.NewUser(tx)
	user, err := userRepo.FindUserByID(ut.UserId)

	if err = mergeData(req, user); err != nil {
		return err
	}

	if err = userRepo.Update(user, strings.Split(allowUpdateField, ",")...); err != nil {
		return err
	}

	return nil
}

func mergeData(req idl.UserUpdateReq, user *idl.User) error {
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}

	if req.Gender != nil {
		user.Gender = *req.Gender
	}

	if req.Birthday != nil {
		user.Birthday = req.Birthday
	}

	if req.ShortDesc != nil {
		user.ShortDesc = *req.ShortDesc
	}

	if req.Avatar != "" {
		if _uuid, err := uuid.Parse(req.Avatar); err != nil {
			return err
		} else {
			user.Avatar = _uuid
		}
	}

	if req.Cover != "" {
		if _uuid, err := uuid.Parse(req.Cover); err != nil {
			return err
		} else {
			user.Cover = _uuid
		}
	}

	return nil
}
