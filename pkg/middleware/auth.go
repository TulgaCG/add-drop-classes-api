package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

func AuthMiddleware(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		token := c.Request.Header.Get(common.TokenHeaderKey)

		user, err := db.GetUserByUsername(context.Background(), username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.Response{
				Error: "login user not found",
			})
			return
		}

		if user.Token.String != token {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, common.Response{
				Error: "login required",
			})
			return
		}

		if time.Since(user.TokenExpireAt.Time) > 0 {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, common.Response{
				Error: "token expired",
			})
			return
		}

		c.Next()
	}
}
