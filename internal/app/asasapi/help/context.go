package help

import (
	"context"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/idl"
	"github.com/gin-gonic/gin"
)

const (
	UserAccessTokenKey = "UserAccessToken"
)

func SetContext(ctx *gin.Context, key string, val interface{}) *gin.Context {
	ctx.Set(key, val)
	return ctx
}

func SetUserToken(ctx *gin.Context, token *idl.Token) *gin.Context {
	return SetContext(ctx, UserAccessTokenKey, token)
}

func GetUserToken(ctx context.Context) *idl.Token {
	if val := ctx.Value(UserAccessTokenKey); val != nil {
		if token, ok := val.(*idl.Token); ok {
			return token
		}
	}
	return nil
}
