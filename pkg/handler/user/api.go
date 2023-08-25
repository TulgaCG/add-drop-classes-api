package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/domain/user"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func List(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	rows, err := user.ListUsers(c, ac.DB)
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to list users")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(rows))
}

func Get(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		ac.Log.Error(response.ErrFailedToFindRolesInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles from context")))
		return
	}

	roles, ok := r.([]types.Role)
	if !ok {
		ac.Log.Error("failed type assertion roles")
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles")))
		return
	}

	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			ac.Log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		self, err := user.GetUserByUsername(c, ac.DB, username)
		if err != nil {
			ac.Log.Error(err.Error())
			c.JSON(http.StatusInternalServerError, response.WithError(response.ErrContentNotFound))
			return
		}

		c.JSON(http.StatusOK, response.WithData(self))
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	row, err := user.GetUser(c, ac.DB, types.UserID(id))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get user")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func Update(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	var req UpdateUserRequest
	if err := c.BindJSON(&req); err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	r, exists := c.Get(common.RolesCtxKey)
	if !exists {
		ac.Log.Error(response.ErrFailedToFindRolesInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles from context")))
		return
	}

	roles, ok := r.([]types.Role)
	if !ok {
		ac.Log.Error("failed type assertion roles")
		c.JSON(http.StatusInternalServerError, response.WithError(errors.New("failed to get roles")))
		return
	}

	if !slices.Contains(roles, types.RoleAdmin) && !slices.Contains(roles, types.RoleTeacher) {
		username := c.Request.Header.Get(common.UsernameHeaderKey)
		if username == "" {
			ac.Log.Error("no username header")
			c.JSON(http.StatusUnauthorized, response.WithError(response.ErrFailedToAuthenticate))
			return
		}

		req.Username = username
	}

	u, err := user.GetUserByUsername(c, ac.DB, req.Username)
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(errors.New("failed to get user by username")))
		return
	}

	row, err := user.UpdateUser(c, ac.DB, gendb.UpdateUserParams{
		ID:       u.ID,
		Username: req.NewUsername,
		Password: req.NewPassword,
	})
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(errors.New("failed to update user")))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}
