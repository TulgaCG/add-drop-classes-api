package role

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

func Create(c *gin.Context) {
	var role gendb.Role
	err := c.Bind(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "bad request",
		})
		return
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to get db",
		})
		return
	}

	addedRole, err := db.CreateRole(context.Background(), role.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to create role",
		})
	}

	c.JSON(http.StatusOK, common.Response{
		Data: addedRole,
	})
}

func Delete(c *gin.Context) {
	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to get db",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("role"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "failed to get role id",
		})
		return
	}

	err = db.DeleteRoleByID(context.Background(), int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to delete role",
		})
	}
}

func AddToUser(_ *gin.Context) {
}

func RemoveFromUser(_ *gin.Context) {
}
