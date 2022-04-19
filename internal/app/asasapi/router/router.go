package router

import (
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/handler"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/middlewares"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/service"
	"github.com/A-SoulFan/asasfans-api/internal/pkg/httpserver"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Provide() fx.Option {
	return fx.Provide(InitRouters)
}

func InitRouters(
	bvService *service.BilbilVideo,
	authService *service.Auth,
	errMiddlewares *middlewares.ErrorInterceptor,
	sessionMiddlewares *middlewares.Session,
) httpserver.InitRouters {
	return func(r *gin.Engine) {
		// http 异常处理
		r.Use(errMiddlewares.Handler)

		// 视频搜素
		r.GET("/v2/asoul-video-interface/advanced-search", handler.BilibiliVideoSearch(bvService))

		// Auth相关
		authApi := r.Group("/v2/auth")
		{
			authApi.POST("/email/send-verifycode", handler.SendEmailVerifyCode(authService))
			authApi.POST("/email/register", handler.EmailRegister(authService))
			authApi.POST("/email/password-signin", handler.EmailPasswordSignIn(authService))
			authApi.POST("/email/password-reset", handler.EmailRestPassword(authService))
		}

		// User相关
		usersApi := r.Group("/v2/users").Use(sessionMiddlewares.Handler)
		{
			usersApi.GET("/info")
		}
	}
}
