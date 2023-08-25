package role

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/appctx"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func AddToUserHandler(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	var req AddRoleRequest
	if err := c.BindJSON(&req); err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	row, err := addRoleToUser(c, ac.DB, req.UserID, req.RoleID)
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func RemoveFromUserHandler(c *gin.Context) {
	ac, ok := c.MustGet(common.AppCtx).(*appctx.Context)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindAppCtxInCtx))
		return
	}

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	rid, err := strconv.Atoi(c.Param("rid"))
	if err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	if err := removeRoleFromUser(c, ac.DB, types.UserID(uid), types.RoleID(rid)); err != nil {
		ac.Log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	c.JSON(http.StatusOK, response.WithData(""))
}
