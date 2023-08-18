package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

func DBMiddleware(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.DatabaseCtxKey, db)
		c.Next()
	}
}
