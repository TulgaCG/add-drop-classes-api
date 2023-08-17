package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/key"
)

func DBMiddleware(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key.DatabaseCtxKey, db)
		c.Next()
	}
}
