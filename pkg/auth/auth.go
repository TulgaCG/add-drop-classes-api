package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/middleware"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

const (
	tokenLength = 64
	validTime   = 1 * time.Hour
)

type UserParams struct {
	TokenExpireAt time.Time      `json:"tokenExpireAt"`
	Username      string         `json:"username"`
	Password      string         `json:"password"`
	Token         sql.NullString `json:"token"`
	ID            types.UserID   `json:"id"`
}

func GetByUsername(c *gin.Context) {
	var userToGet UserParams
	if err := c.Bind(&userToGet); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	db, ok := c.MustGet(middleware.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get db",
		})
		return
	}

	user, err := db.GetUserByUsername(context.Background(), userToGet.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context) {
	var userToGet UserParams
	if err := c.Bind(&userToGet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	db, ok := c.MustGet(middleware.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	user, err := db.GetUserByUsername(context.Background(), userToGet.Username)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	if user.Password != userToGet.Password {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if time.Since(user.TokenExpireAt.Time.Add(validTime)) > 0 {
		token, err := db.AddToken(context.Background(), gendb.AddTokenParams{
			ID:    user.ID,
			Token: sql.NullString{String: createRandomToken(), Valid: true},
		})
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":    token,
			"username": user.Username,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"token":    user.Token,
			"username": user.Username,
			"message":  "logged in",
		})
	}
}

func Logout(c *gin.Context) {
	var userToGet UserParams
	if err := c.Bind(&userToGet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	db, ok := c.MustGet(middleware.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	user, err := db.GetUserByUsername(context.Background(), userToGet.Username)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	if user.Token != userToGet.Token {
		c.Status(http.StatusNotAcceptable)
		return
	}

	_, err = db.ExpireToken(context.Background(), gendb.ExpireTokenParams{
		TokenExpireAt: sql.NullTime{
			Time:  time.Now().Add(-validTime),
			Valid: true,
		},
		ID: user.ID,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    "",
		"username": "",
		"message":  "logged out",
	})
}

func createRandomToken() string {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
