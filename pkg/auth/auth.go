package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/TulgaCG/add-drop-classes-api/pkg/common"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

type UserParams struct {
	TokenExpireAt time.Time      `json:"tokenExpireAt"`
	Username      string         `json:"username"`
	Password      string         `json:"password"`
	Token         sql.NullString `json:"token"`
	ID            types.UserID   `json:"id"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func Login(c *gin.Context) {
	var userToGet UserParams
	if err := c.Bind(&userToGet); err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "bad request",
		})
	}

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to get db",
		})
		return
	}

	user, err := db.GetUserByUsername(context.Background(), userToGet.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Error: "failed to get user by username",
		})
		return
	}

	if user.Password != userToGet.Password {
		c.JSON(http.StatusNotAcceptable, common.Response{
			Error: "wrong password",
		})
		return
	}

	if time.Since(user.TokenExpireAt.Time) > 0 {
		generatedToken, err := createRandomToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Response{
				Error: "failed to generate token",
			})
		}

		token, err := db.UpdateToken(context.Background(), gendb.UpdateTokenParams{
			ID:    user.ID,
			Token: sql.NullString{String: generatedToken, Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Response{
				Error: "failed to update token",
			})
			return
		}

		_, err = db.UpdateExpirationToken(context.Background(), gendb.UpdateExpirationTokenParams{
			ID:            user.ID,
			TokenExpireAt: sql.NullTime{Time: time.Now().Add(common.ValidTime), Valid: true},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Response{
				Error: "failed to update token expiration date",
			})
			return
		}

		c.JSON(http.StatusOK, common.Response{
			Data: LoginResponse{
				Username: user.Username,
				Token:    token.String,
			},
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			Data: LoginResponse{
				Username: user.Username,
				Token:    user.Token.String,
			},
		})
	}
}

func Logout(c *gin.Context) {
	username := c.Request.Header.Get(common.UsernameHeaderKey)
	token := c.Request.Header.Get(common.TokenHeaderKey)

	db, ok := c.MustGet(common.DatabaseCtxKey).(*gendb.Queries)
	if !ok {
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to get db",
		})
		return
	}

	user, err := db.GetUserByUsername(context.Background(), username)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, common.Response{
			Error: "username or database not found",
		})
		return
	}

	if user.Token.String != token {
		c.JSON(http.StatusNotAcceptable, common.Response{
			Error: "not logged in",
		})
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
		c.JSON(http.StatusInternalServerError, common.Response{
			Error: "failed to update token expiration date",
		})
	}

	c.JSON(http.StatusOK, common.Response{
		Data: username,
	})
}

func createRandomToken() (string, error) {
	b := make([]byte, common.TokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
