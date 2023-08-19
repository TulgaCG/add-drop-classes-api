package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

const (
	tokenLen      = 64
	tokenDuration = 1 * time.Hour
)

func login(ctx context.Context, db *gendb.Queries, username, password string) (gendb.UpdateTokenRow, error) {
	u, err := getUserCredentialsWithUsername(ctx, db, username)
	if err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("failed to get user with username including credentials: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("invalid username or password")
	}

	token, err := createRandomToken()
	if err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("failed to create token: %w", err)
	}

	row, err := db.UpdateToken(ctx, gendb.UpdateTokenParams{
		ID: u.ID,
		Token: sql.NullString{
			String: token,
			Valid:  true,
		},
		TokenExpireAt: sql.NullTime{
			Time:  time.Now().Add(tokenDuration),
			Valid: true,
		},
	})
	if err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("failed to update token: %w", err)
	}

	return row, nil
}

func logout(ctx context.Context, db *gendb.Queries, username string) error {
	user, err := db.GetUserCredentialsWithUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user by username: %w", err)
	}

	if err := db.UpdateTokenExpirationDate(ctx, gendb.UpdateTokenExpirationDateParams{
		TokenExpireAt: sql.NullTime{
			Time:  time.Now(),
			Valid: false,
		},
		ID: user.ID,
	}); err != nil {
		return fmt.Errorf("failed to update token expiration date: %w", err)
	}

	return nil
}

func getUserCredentialsWithUsername(ctx context.Context, db *gendb.Queries, username string) (gendb.User, error) {
	row, err := db.GetUserCredentialsWithUsername(ctx, username)
	if err != nil {
		return gendb.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return row, nil
}

func createRandomToken() (string, error) {
	b := make([]byte, tokenLen)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return hex.EncodeToString(b), nil
}
