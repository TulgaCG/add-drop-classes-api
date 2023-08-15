package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gen/db"
)

const DatabaseCtxKey = "db"

func DBMiddleware(db *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(DatabaseCtxKey, db)
		c.Next()
	}
}
