package handler

import (
	"asasfans/internal/app/asasapi/help"
	"asasfans/internal/app/asasapi/idl"
	"asasfans/internal/app/asasapi/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BilbilVideoSearch(s *service.BilbilVideo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req idl.BilbilVideoSearchReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusOK, help.FailureJson(404, err.Error(), nil))
			return
		}

		if resp, err := s.Search(ctx, req); err != nil {
			ctx.JSON(http.StatusInternalServerError, help.FailureJson(404, err.Error(), nil))
			return
		} else {
			ctx.JSON(http.StatusOK, help.SuccessJson(resp))
		}
	}
}
