package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
)

// AppContext adds appctx to the gin.Context
// See appctx.Context for what it includes.
func AppContext(appctx *appctx.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.AppCtx, appctx)
		c.Next()
	}
}
