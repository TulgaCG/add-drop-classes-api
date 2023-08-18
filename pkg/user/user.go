package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func Post(c *gin.Context) {
	var req gendb.CreateUserParams
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

	newUser, err := db.CreateUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Data: newUser,
	})
}

func Get(c *gin.Context) {
	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to get db",
		})
		return
	}

	users, err := db.ListUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusNoContent, common.Response{
			Error: "failed to list users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetByID(c *gin.Context) {
	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to get db",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "id must be integer",
		})
		return
	}

	u, err := db.GetUser(context.Background(), types.UserID(id))
	if err != nil {
		c.JSON(http.StatusNoContent, common.Response{
			Error: "failed to get user by id",
		})
		return
	}

	c.JSON(http.StatusOK, u)
}

func Update(c *gin.Context) {
	var req gendb.UpdateUserParams
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

	u, err := db.UpdateUser(context.Background(), gendb.UpdateUserParams{
		Username: req.Username,
		Password: req.Password,
		ID:       req.ID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, u)
}

func Delete(c *gin.Context) {
	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to get db",
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "failed to get id",
		})
		return
	}

	rows, err := db.DeleteUser(context.Background(), types.UserID(id))
	if err != nil {
		c.JSON(http.StatusNoContent, common.Response{
			Error: "failed to delete user",
		})
		return
	}

	if rows <= 0 {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "user to delete not found",
		})
	}

	c.JSON(http.StatusOK, common.Response{})
}
