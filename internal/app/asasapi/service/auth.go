package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/util/password"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/cache"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/smsclient"
	"gorm.io/gorm"
)

const (
	typeSend = "SEND_VERIFY_CODE"
	typeSave = "SAVE_VERIFY_CODE"
)

type Auth struct {
	db            *gorm.DB
	cache         cache.ICache
	pwdHandler    password.PasswordHandler
	codeGenerator smsclient.CodeGenerator
	smsClient     smsclient.Client
}

func NewAuth(
	db *gorm.DB,
	cache cache.ICache,
	passwordHandler password.PasswordHandler,
	codeGenerator smsclient.CodeGenerator,
	smsClient smsclient.Client,
) *Auth {
	return &Auth{
		db:            db,
		cache:         cache,
		pwdHandler:    passwordHandler,
		codeGenerator: codeGenerator,
		smsClient:     smsClient,
	}
}

func (a *Auth) SendEmailVerifyCode(ctx context.Context, req idl.SendEmailVerifyCodeReq) error {
	sendKey, saveCodeKey := buildKey(req.Email, typeSend), buildKey(req.Email, typeSave)

	if _, isset := a.cache.Get(sendKey); isset {
		return apperrors.NewValidationError(SmsSendLimitError, "发送验证码过于频繁")
	}

	verifyCode := a.codeGenerator.Generate(6)
	if err := a.cache.Set(sendKey, time.Now().Unix(), 50*time.Second); err != nil {
		return err
	}

	if err := a.cache.Set(saveCodeKey, verifyCode, 5*time.Minute); err != nil {
		_ = a.cache.Delete(saveCodeKey)
		return err
	}

	if err := a.smsClient.SendHTML(req.Email, verifyCode, "邮箱验证码"); err != nil {
		_ = a.cache.Delete(saveCodeKey)
		return err
	}

	return nil
}

func (a *Auth) EmailRegister(ctx context.Context, req idl.EmailRegisterReq) error {
	saveCodeKey := buildKey(req.Email, typeSave)

	defer a.cache.Delete(saveCodeKey)
	if code, isset := a.cache.Get(saveCodeKey); !isset || code != req.VerifyCode {
		return apperrors.NewValidationError(InvalidVerifyCodeError, "无效的验证码")
	}

	tx := a.db.WithContext(ctx)
	userRep := repository.NewUser(tx)

	if u, err := userRep.FindUserByEmail(req.Email); err != nil {
		return err
	} else if u != nil {
		return apperrors.NewValidationError(RepeatRegisteredError, "该邮箱已被注册")
	}

	user := DefaultNewUser(req.Email)
	user.Password = a.pwdHandler.Hash(req.Password)

	if isNew, err := userRep.Create(user); err != nil {
		return err
	} else if !isNew {
		return apperrors.NewValidationError(RepeatRegisteredError, "该邮箱已被注册")
	}

	return nil
}

func (a *Auth) EmailPasswordSignIn(ctx context.Context, req idl.EmailPasswordSignInReq) (*idl.SignInResp, error) {
	tx := a.db.WithContext(ctx)

	user, err := repository.NewUser(tx).FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, apperrors.NewValidationError(InviteAccount, "用户名或密码错误")
	}

	if !a.pwdHandler.Check(user.Password, req.Password) {
		return nil, apperrors.NewValidationError(InviteAccount, "用户名或密码错误")
	}

	token := buildSessionKey(user)
	if err = repository.NewTokenStorage(tx).Set(token, &idl.Token{UserId: user.Id}); err != nil {
		return nil, err
	}

	return &idl.SignInResp{Token: token}, nil
}

func (a *Auth) EmailRestPassword(ctx context.Context, req idl.RestPasswordReq) error {
	tx := a.db.WithContext(ctx)

	user, err := repository.NewUser(tx).FindUserByEmail(req.Email)
	if err != nil {
		return err
	} else if user == nil {
		return apperrors.NewValidationError(InviteAccount, "无效的账户")
	}

	err = tx.Transaction(func(_tx *gorm.DB) error {
		if err := repository.NewUser(_tx).Update(user, "password"); err != nil {
			return err
		}

		if err := repository.NewTokenStorage(_tx).ClearTokens(&idl.Token{UserId: user.Id}); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func DefaultNewUser(email string) *idl.User {
	return &idl.User{
		Nickname:  "新用户",
		Email:     email,
		Gender:    idl.UserGenderDefault,
		Birthday:  nil,
		ShortDesc: "最初的个人简介，送给一个魂们~",
		Status:    idl.UserStatusNormal,
	}
}

func buildKey(k, t string) string {
	return t + "_" + k
}

func buildSessionKey(user *idl.User) string {
	hx := sha1.New()
	hx.Write([]byte(fmt.Sprintf("%d-%d-%s", time.Now().UnixNano(), user.Id, user.Email)))

	return hex.EncodeToString(hx.Sum(nil))
}
