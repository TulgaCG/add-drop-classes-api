package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
)

const tokenHeaderKey = "Token"

func AuthView(db *gendb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
			return
		}

		username, token, err := checkCookiesAndHeaders(c)
		if err != nil {
			ac.Log.Error(err.Error())
			c.Redirect(http.StatusFound, "/login")
		}

		u, err := db.GetUserCredentialsWithUsername(c, username)
		if err != nil {
			ac.Log.Error(fmt.Sprintf("failed to get user credentials with username %s: %s", username, err.Error()))
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if u.Token.String != token {
			ac.Log.Error("failed to match user token with the token")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		if time.Since(u.TokenExpireAt.Time) > 0 {
			ac.Log.Error("token expired")
			c.Redirect(http.StatusFound, "/login")
			return
		}

		c.Next()
	}
}

func Auth(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
		if !ok {
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
			return
		}

		username, token, err := checkCookiesAndHeaders(c)
		if err != nil {
			ac.Log.Error(err.Error())
			c.Redirect(http.StatusFound, "/login")
		}

		u, err := db.GetUserCredentialsWithUsername(c, username)
		if err != nil {
			ac.Log.Error(fmt.Sprintf("no database in gin context %s", username))
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		if u.Token.String != token {
			ac.Log.Error("failed to match user token with the header token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		if time.Since(u.TokenExpireAt.Time) > 0 {
			ac.Log.Error("token expired")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.WithError(fmt.Errorf("token expired")))
			return
		}

		c.Next()
	}
}

func checkCookiesAndHeaders(c *gin.Context) (string, string, error) {
	username := c.Request.Header.Get(common.UsernameHeaderKey)
	token := c.Request.Header.Get(tokenHeaderKey)
	if username == "" || token == "" {
		var err error
		username, err = c.Cookie("username")
		if err != nil {
			return "", "", fmt.Errorf("failed to find username")
		}

		token, err = c.Cookie("token")
		if err != nil {
			return "", "", fmt.Errorf("failed to find token")
		}
	}

	return username, token, nil
}
