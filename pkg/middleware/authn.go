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

func AuthenticationUI(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
			return
		}

		username := c.Request.Header.Get(common.UsernameHeaderKey)
		token := c.Request.Header.Get(tokenHeaderKey)

		if username == "" {
			username, err = c.Cookie("username")
			if err != nil {
				log.Error("no username found")
				c.Redirect(http.StatusFound, "/login")
				return
			}

			token, err = c.Cookie("token")
			if err != nil {
				log.Error("no username found")
				c.Redirect(http.StatusFound, "/login")
				return
			}
		}

		u, err := db.GetUserCredentialsWithUsername(c, username)
		if err != nil {
			log.Error(fmt.Sprintf("no database in gin context %s", username))
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if u.Token.String != token {
			log.Error("failed to match user token with the header token")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if time.Since(u.TokenExpireAt.Time) > 0 {
			log.Error("token expired")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.Next()
	}
}

func Authentication(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
			return
		}

		username := c.Request.Header.Get(common.UsernameHeaderKey)
		token := c.Request.Header.Get(tokenHeaderKey)

		if username == "" {
			username, err = c.Cookie("username")
			if err != nil {
				log.Error("no username found")
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
				return
			}

			token, err = c.Cookie("token")
			if err != nil {
				log.Error("no username found")
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
				return
			}
		}

		u, err := db.GetUserCredentialsWithUsername(c, username)
		if err != nil {
			log.Error(fmt.Sprintf("no database in gin context %s", username))
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