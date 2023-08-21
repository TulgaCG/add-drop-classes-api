package role

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/server/response"
)

func CreateRoleHandler(c *gin.Context) {
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

	var role CreateRoleRequest
	if err := c.BindJSON(&role); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidRequestFormat))
		return
	}

	row, err := createRole(c, db, role.RoleName)
	if err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, response.WithError(err))
		return
	}

	c.JSON(http.StatusOK, response.WithData(row))
}

func DeleteRoleHandler(c *gin.Context) {
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

	role := c.Param("role")
	if err := deleteRole(c, db, role); err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusBadRequest, response.WithError(response.ErrInvalidParamIDFormat))
		return
	}

	c.JSON(http.StatusOK, response.WithData(""))
}
