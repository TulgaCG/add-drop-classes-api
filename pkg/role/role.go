package role

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

type ToUser struct {
	Username string       `json:"username"`
	Rolename string       `json:"rolename"`
	UserID   types.UserID `json:"userID"`
	RoleID   types.RoleID `json:"roleID"`
}

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

	err = db.DeleteRoleByID(context.Background(), types.RoleID(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to delete role",
		})
	}
}

func AddToUser(c *gin.Context) {
	var req ToUser
	err := c.Bind(&req)
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

	var userToAdd gendb.User
	var roleToAdd gendb.Role

	if req.Username != "" {
		userToAdd, err = db.GetUserByUsername(context.Background(), req.Username)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Response{
				Error: "username not found",
			})
			return
		}
	} else {
		userToAdd, err = db.GetUser(context.Background(), req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Response{
				Error: "user not found",
			})
			return
		}
	}

	if req.Rolename != "" {
		roleToAdd, err = db.GetRoleByName(context.Background(), req.Rolename)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Response{
				Error: "rolename not found",
			})
			return
		}
	} else {
		roleToAdd, err = db.GetRole(context.Background(), req.RoleID)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Response{
				Error: "role not found",
			})
			return
		}
	}

	err = db.CreateUserRole(context.Background(), gendb.CreateUserRoleParams{
		UserID: userToAdd.ID,
		RoleID: roleToAdd.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to add role to the user",
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Data: struct {
			UserID int64 `json:"userID"`
			RoleID int64 `json:"roleID"`
		}{
			UserID: int64(userToAdd.ID),
			RoleID: int64(roleToAdd.ID),
		},
	})
}

func RemoveFromUser(c *gin.Context) {
	userid, err := strconv.Atoi(c.Param("user"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "failed to get user id",
		})
		return
	}

	roleid, err := strconv.Atoi(c.Param("role"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "failed to get role id",
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

	rows, err := db.DeleteUserRole(context.Background(), gendb.DeleteUserRoleParams{
		UserID: types.UserID(userid),
		RoleID: types.RoleID(roleid),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to delete role from user",
		})
		return
	}

	if rows <= 0 {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "no data found to delete",
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{})
}
