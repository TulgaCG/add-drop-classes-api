package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
)

const (
	cookieLifeTime = 7200
)

func LoginHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	var req LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	row, err := login(c, db, req.Username, req.Password)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
		return
	}

	c.SetCookie("username", row.Username, cookieLifeTime, "/", "127.0.0.1", false, true)
	c.SetCookie("token", row.Token.String, cookieLifeTime, "/", "127.0.0.1", false, true)

	c.JSON(http.StatusOK, response.WithData(row))
}

func LogoutHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	username := c.Request.Header.Get(common.UsernameHeaderKey)
	if username == "" {
		log.Error(response.ErrFailedToAuthenticate.Error())
		c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	if err := logout(c, db, username); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	c.JSON(http.StatusOK, response.WithData(nil))
}
