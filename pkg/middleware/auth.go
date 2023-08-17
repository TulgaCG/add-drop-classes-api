package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/key"
)

func AuthMiddleware(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.Header.Get(key.UsernameHeaderKey)
		token := c.Request.Header.Get(key.TokenHeaderKey)

		user, err := db.GetUserByUsername(context.Background(), username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "login user not found",
			})
			return
		}

		if user.Token.String != token {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"error": "login required",
			})
			return
		}

		if time.Since(user.TokenExpireAt.Time) > 0 {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"error": "token expired",
			})
			return
		}

		c.Next()
	}
}
