package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/domain/auth"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
)

const (
	cookieLifeTime = 7200
)

func Login(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	var req LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	row, err := auth.Login(c, ac.DB, req.Username, req.Password)
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
		return
	}

	c.SetCookie("username", row.Username, cookieLifeTime, "/", "127.0.0.1", false, true)
	c.SetCookie("token", row.Token.String, cookieLifeTime, "/", "127.0.0.1", false, true)

	c.JSON(http.StatusOK, response.WithData(row))
}

func Logout(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	username := c.Request.Header.Get(common.UsernameHeaderKey)
	if username == "" {
		ac.Log.Error(response.ErrFailedToAuthenticate.Error())
		c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
		return
	}

	if err := auth.Logout(c, ac.DB, username); err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusUnauthorized, response.WithError(errors.New("failed to logout")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(nil))
}

func Register(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	row, err := auth.Register(c, ac.DB, req.Username, req.Password)
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to register user, perhaps user already exists")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}
