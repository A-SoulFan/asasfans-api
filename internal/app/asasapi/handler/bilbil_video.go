package handler

import (
	"net/http"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/help"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/service"

	"github.com/gin-gonic/gin"
)

func BilibiliVideoSearch(s *service.BilbilVideo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req idl.BilibiliVideoSearchReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			_ = ctx.Error(apperrors.NewValidationError(400, err.Error()).Wrap(err))
			return
		}

		if resp, err := s.Search(ctx, req); err != nil {
			_ = ctx.Error(err)
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(resp))
		}
	}
}
