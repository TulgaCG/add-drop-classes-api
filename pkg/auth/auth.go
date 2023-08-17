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
	"github.com/TulgaCG/add-drop-classes-api/pkg/key"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

type UserParams struct {
	TokenExpireAt time.Time      `json:"tokenExpireAt"`
	Username      string         `json:"username"`
	Password      string         `json:"password"`
	Token         sql.NullString `json:"token"`
	ID            types.UserID   `json:"id"`
}

func Login(c *gin.Context) {
	var userToGet UserParams
	if err := c.Bind(&userToGet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
	}

	db, ok := c.MustGet(key.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	user, err := db.GetUserByUsername(context.Background(), userToGet.Username)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if user.Password != userToGet.Password {
		c.Status(http.StatusNotAcceptable)
		return
	}

	if time.Since(user.TokenExpireAt.Time) > 0 {
		token, err := db.UpdateToken(context.Background(), gendb.UpdateTokenParams{
			ID:    user.ID,
			Token: sql.NullString{String: createRandomToken(), Valid: true},
		})
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		_, err = db.UpdateExpirationToken(context.Background(), gendb.UpdateExpirationTokenParams{
			ID:            user.ID,
			TokenExpireAt: sql.NullTime{Time: time.Now().Add(key.ValidTime), Valid: true},
		})
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token":    token,
			"username": user.Username,
			"message":  "logged in",
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
	username := c.Request.Header.Get(key.UsernameHeaderKey)
	token := c.Request.Header.Get(key.TokenHeaderKey)

	db, ok := c.MustGet(key.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.Status(http.StatusInternalServerError)
		return
	}

	user, err := db.GetUserByUsername(context.Background(), username)
	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	if user.Token.String != token {
		c.Status(http.StatusNotAcceptable)
		return
	}

	_, err = db.UpdateExpirationToken(context.Background(), gendb.UpdateExpirationTokenParams{
		TokenExpireAt: sql.NullTime{
			Time:  time.Now(),
			Valid: false,
		},
		ID: user.ID,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"message":  "logged out",
	})
}

func createRandomToken() string {
	b := make([]byte, key.TokenLength)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
