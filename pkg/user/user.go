package user

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gen/db"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
)

func Post(c *gin.Context) {
	// Get request
	var req db.CreateUserParams
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	// Get query
	query, ok := c.MustGet(middleware.DatabaseCtxKey).(*db.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get db",
		})
		return
	}

	// Create user
	newUser, err := query.CreateUser(context.Background(), db.CreateUserParams{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": newUser,
	})
}

func Get(c *gin.Context) {
	// Get query
	query, ok := c.MustGet(middleware.DatabaseCtxKey).(*db.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get db",
		})
		return
	}

	// Get list
	users, err := query.ListUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "failed to list users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func GetByID(c *gin.Context) {
	// Get query
	query, ok := c.MustGet(middleware.DatabaseCtxKey).(*db.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get db",
		})
		return
	}

	// Get id from url
	var id sql.NullInt64
	err := id.Scan(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id must be integer",
		})
		return
	}

	// Get user
	u, err := query.GetUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "failed to get user by id",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

func GetUserByUsername(c *gin.Context) {
	// Get query
	query, ok := c.MustGet(middleware.DatabaseCtxKey).(*db.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get db",
		})
		return
	}

	// Get username from url
	var username sql.NullString
	err := username.Scan(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get username",
		})
		return
	}

	// Get user
	u, err := query.GetUserByUsername(context.Background(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get user by username",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

func Update(c *gin.Context) {
	// Get request
	var req db.UpdateUserParams
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	// Get query
	query, ok := c.MustGet(middleware.DatabaseCtxKey).(*db.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get db",
		})
		return
	}

	// Update user
	u, err := query.UpdateUser(context.Background(), db.UpdateUserParams{
		Username: req.Username,
		Password: req.Password,
		ID:       req.ID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})
}

func Delete(c *gin.Context) {
	// Get query
	query, ok := c.MustGet(middleware.DatabaseCtxKey).(*db.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get db",
		})
		return
	}

	// Get id from url
	var id sql.NullInt64
	err := id.Scan(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get id",
		})
		return
	}

	// Delete user
	err = query.DeleteUser(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{
			"error": "failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": "user deleted successfully",
	})
}
