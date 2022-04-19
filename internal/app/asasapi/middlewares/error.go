package middlewares

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/apperrors"
	"github.com/A-SoulFan/asasfans-api/internal/app/asasapi/help"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type ErrorInterceptor struct {
	logger *zap.Logger
}

func NewErrorInterceptor(logger *zap.Logger) *ErrorInterceptor {
	return &ErrorInterceptor{logger: logger}
}

func (e *ErrorInterceptor) Handler(ctx *gin.Context) {
	defer func() {
		if len(ctx.Errors) == 0 {
			ctx.Next()
			return
		}

		headers, _ := json.Marshal(ctx.Request.Header)
		logs := []zap.Field{
			zap.String("request.method", ctx.Request.Method),
			zap.String("request.url", ctx.Request.URL.String()),
			zap.ByteString("request.headers", headers),
			zap.String("errors", ginErrorsToString(ctx.Errors)),
		}
		e.logger.Error("request error:", logs...)

		code, msg := -1, "服务器异常，请稍后再试"
		for i := len(ctx.Errors) - 1; i >= 0; i-- {
			err := ctx.Errors[i]
			if appError, ok := errors.Cause(err.Err).(*apperrors.AppError); ok {
				code = appError.Code
				switch appError.ResponseType {
				case apperrors.ValidationError:
					msg = appError.Message
				case apperrors.AuthError:
					msg = appError.Message
				}
				break
			}
		}

		ctx.JSON(http.StatusOK, help.FailureJson(code, msg, nil))
	}()
	ctx.Next()
}

func ginErrorsToString(errs []*gin.Error) string {
	if len(errs) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, msg := range errs {
		_, _ = fmt.Fprintf(&buffer, "Error #%02d: %+v\n", i+1, msg.Err)
		if msg.Meta != nil {
			_, _ = fmt.Fprintf(&buffer, "Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
