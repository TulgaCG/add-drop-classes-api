package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
)

func LogMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.LogCtxKey, log)
		c.Next()
	}
}
