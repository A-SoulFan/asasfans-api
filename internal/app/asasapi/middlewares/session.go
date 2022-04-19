package middlewares

import (
	"net/http"
	"strings"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/help"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Session struct {
	db *gorm.DB
}

func NewSession(db *gorm.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Handler(ctx *gin.Context) {
	token := getToken(ctx.Request)
	if token == "" {
		_ = ctx.Error(apperrors.NewAuthError(401, "not found token"))
		ctx.Abort()
		return
	}

	tokenStorage := repository.NewTokenStorage(s.db.WithContext(ctx))
	ut, err := tokenStorage.Get(token)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Abort()
		return
	}

	help.SetUserToken(ctx, ut)
	ctx.Next()
}

func getToken(r *http.Request) string {
	authorization := r.Header.Get("authorization")
	if len(authorization) >= 20 {
		token := strings.Split(authorization, " ")
		if len(token) == 2 && len(token[1]) >= 20 {
			return token[1]
		}
	}

	return ""
}
