package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

const (
	defaultRole   = types.RoleID(3)
	tokenCharset  = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	tokenLen      = 64
	tokenDuration = 1 * time.Hour
)

func Login(ctx context.Context, db *database.DB, username, password string) (gendb.UpdateTokenRow, error) {
	u, err := getUserCredentialsWithUsername(ctx, db, username)
	if err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("failed to get user with username including credentials: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("invalid username or password")
	}

	token, err := createRandomToken(tokenLen)
	if err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("failed to create token: %w", err)
	}

	row, err := db.UpdateToken(ctx, gendb.UpdateTokenParams{
		ID: u.ID,
		Token: pgtype.Text{
			String: token,
			Valid:  true,
		},
		TokenExpireAt: pgtype.Timestamp{
			Time:  time.Now().Add(tokenDuration),
			Valid: true,
		},
	})
	if err != nil {
		return gendb.UpdateTokenRow{}, fmt.Errorf("failed to update token: %w", err)
	}

	return row, nil
}

func Logout(ctx context.Context, db *database.DB, username string) error {
	user, err := db.GetUserCredentialsWithUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to get user by username: %w", err)
	}

	if err := db.UpdateTokenExpirationDate(ctx, gendb.UpdateTokenExpirationDateParams{
		TokenExpireAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: false,
		},
		ID: user.ID,
	}); err != nil {
		return fmt.Errorf("failed to update token expiration date: %w", err)
	}

	return nil
}

func Register(ctx context.Context, db *database.DB, username, password string) (gendb.CreateUserRow, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return gendb.CreateUserRow{}, fmt.Errorf("failed to generate hashed password: %w", err)
	}

	var row gendb.CreateUserRow
	if err := db.WithTx(ctx, func(dbTx *gendb.Queries) error {
		var err error
		if row, err = dbTx.CreateUser(ctx, gendb.CreateUserParams{
			Username: username,
			Password: string(hashedPassword),
		}); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		if _, err := dbTx.AddRoleToUser(ctx, gendb.AddRoleToUserParams{
			UserID: row.ID,
			RoleID: defaultRole,
		}); err != nil {
			return fmt.Errorf("failed to add role to the user: %w", err)
		}

		return nil
	}); err != nil {
		return gendb.CreateUserRow{}, fmt.Errorf("failed to create or add role to the user: %w", err)
	}

	return row, nil
}

func getUserCredentialsWithUsername(ctx context.Context, db *database.DB, username string) (gendb.User, error) {
	row, err := db.GetUserCredentialsWithUsername(ctx, username)
	if err != nil {
		return gendb.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return row, nil
}

func createRandomToken(length int) (string, error) {
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(tokenCharset))))
		if err != nil {
			return "", err
		}
		ret[i] = tokenCharset[num.Int64()]
	}

	return string(ret), nil
}
