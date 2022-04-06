package router

import (
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/handler"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/service"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/httpserver"

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
