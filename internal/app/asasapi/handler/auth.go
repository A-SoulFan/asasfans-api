package handler

import (
	"net/http"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/help"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/service"
	"github.com/gin-gonic/gin"
)

func SendEmailVerifyCode(a *service.Auth) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req idl.SendEmailVerifyCodeReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			_ = ctx.Error(apperrors.NewValidationError(400, err.Error()).Wrap(err))
			return
		}

		if err := a.SendEmailVerifyCode(ctx, req); err != nil {
			_ = ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(nil))
		}
	}
}

func EmailRegister(a *service.Auth) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req idl.EmailRegisterReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			_ = ctx.Error(apperrors.NewValidationError(400, err.Error()).Wrap(err))
			return
		}

		if err := a.EmailRegister(ctx, req); err != nil {
			_ = ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(nil))
		}
	}
}

func EmailPasswordSignIn(a *service.Auth) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req idl.EmailPasswordSignInReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			_ = ctx.Error(apperrors.NewValidationError(400, err.Error()).Wrap(err))
			return
		}

		if resp, err := a.EmailPasswordSignIn(ctx, req); err != nil {
			_ = ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(resp))
		}
	}
}

func EmailRestPassword(a *service.Auth) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req idl.RestPasswordReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			_ = ctx.Error(apperrors.NewValidationError(400, err.Error()).Wrap(err))
			return
		}

		if err := a.EmailRestPassword(ctx, req); err != nil {
			_ = ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(nil))
		}
	}
}
