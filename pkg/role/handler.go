package role

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func AddToUserHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	v, ok := c.MustGet(common.ValidatorCtxKey).(*validator.Validate)
	if !ok {
		log.Error(response.ErrFailedToFindValidatorInCtx.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrFailedToFindValidatorInCtx))
		return
	}

	var req AddRoleRequest
	if err := c.BindJSON(&req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	if err := v.Struct(req); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(fmt.Errorf("failed validation")))
		return
	}

	row, err := addRoleToUser(c, db, req.UserID, req.RoleID)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func RemoveFromUserHandler(c *gin.Context) {
	log, ok := c.MustGet(common.LogCtxKey).(*slog.Logger)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindLoggerInCtx))
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		log.Error(response.ErrFailedToFindDBInCtx.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(response.ErrFailedToFindDBInCtx))
		return
	}

	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	rid, err := strconv.Atoi(c.Param("rid"))
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	if err := removeRoleFromUser(c, db, types.UserID(uid), types.RoleID(rid)); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	c.JSON(http.StatusOK, response.WithData(""))
}
