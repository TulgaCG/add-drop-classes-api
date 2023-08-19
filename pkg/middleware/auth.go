package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
)

const tokenHeaderKey = "Token"

func AuthMiddleware(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
			return
		}

		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			log.Error("no username header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		token := c.Request.Header.Get(tokenHeaderKey)
		if token == "" {
			log.Error("no token header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		u, err := db.GetUserCredentialsWithUsername(c, username)
		if err != nil {
			log.Error("no database in gin context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		if u.Token.String != token {
			log.Error("failed to match user token with the header token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		if time.Since(u.TokenExpireAt.Time) > 0 {
			log.Error("token expired")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(fmt.Errorf("token expired")))
			return
		}

		c.Next()
	}
}
