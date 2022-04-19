package handler

import (
	"net/http"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/help"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/service"
	"github.com/gin-gonic/gin"
)

func UserInfo(u *service.User) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if resp, err := u.Info(ctx); err != nil {
			_ = ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(resp))
		}
	}
}

func UserUpdate(u *service.User) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req idl.UserUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			_ = ctx.Error(apperrors.NewValidationError(400, err.Error()).Wrap(err))
			return
		}

		if err := u.Update(ctx, req); err != nil {
			_ = ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(nil))
		}
	}
}
