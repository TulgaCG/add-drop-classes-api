package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

const DatabaseCtxKey = "db"

func DBMiddleware(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DatabaseCtxKey, db)
		c.Next()
	}
}
