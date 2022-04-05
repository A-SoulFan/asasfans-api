package router

import (
	"asasfans/internal/app/asasapi/handler"
	"asasfans/internal/app/asasapi/service"
	"asasfans/internal/pkg/httpserver"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(InitRouters)
}

func InitRouters(bvService *service.BilbilVideo) httpserver.InitRouters {
	return func(r *gin.Engine) {
		// 视频搜素
		r.GET("/v2/asoul-video-interface/advanced-search", handler.BilbilVideoSearch(bvService))
	}
}
