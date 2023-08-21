package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
)

func TestCreateRandomToken(t *testing.T) {
	token, err := createRandomToken()
	require.NoError(t, err, "createRandomToken failed")

	require.Len(t, token, tokenLen*2, fmt.Sprintf("expected token size: %d, got: %d", tokenLen*2, len(token)))
}

func TestGetUserCredentialsWithUsername(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser1", "testpassword1", false},
		{"wronguser", "testpassword2", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			user, err := getUserCredentialsWithUsername(context.Background(), db, testCase.Username)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(testCase.Password))
				require.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser1", "testpassword1", false},
		{"testuser2", "wrongpassword", true},
		{"wronguser", "testpassword3", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, err := login(context.Background(), db, testCase.Username, testCase.Password)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, row.Token.String, tokenLen*2)
				user, err := db.GetUserCredentialsWithUsername(context.Background(), testCase.Username)
				require.NoError(t, err, "failed to get user by username")
				require.Equal(t, user.Token.String, row.Token.String)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser1", "testpassword1", false},
		{"testuser2", "wrongpassword", false},
		{"wronguser", "testpassword3", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, _ := login(context.Background(), db, testCase.Username, testCase.Password)
			err = logout(context.Background(), db, testCase.Username)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				user, err := getUserCredentialsWithUsername(context.Background(), db, testCase.Username)
				require.NoError(t, err, "failed to get user by username")
				require.Equal(t, user.TokenExpireAt.Valid, false)
				require.Equal(t, user.Token.String, row.Token.String)
			}
		})
	}
}
