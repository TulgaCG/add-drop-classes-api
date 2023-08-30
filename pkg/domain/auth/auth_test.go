package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/domain/user"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

const errCreateMockData = "failed to create mock data"

func TestCreateRandomToken(t *testing.T) {
	token, err := createRandomToken(tokenLen)
	require.NoError(t, err)
	require.Len(t, token, tokenLen)
}

func TestGetUserCredentialsWithUsername(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	require.NoError(t, err, errCreateMockData)
	_, err = db.CreateUser(context.Background(), gendb.CreateUserParams{Username: "testuser", Password: string(hashedPassword)})
	require.NoError(t, err, errCreateMockData)

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser", "testpassword", false},
		{"wronguser", "testpassword2", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			u, err := getUserCredentialsWithUsername(context.Background(), db, testCase.Username)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(testCase.Password))
				require.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	require.NoError(t, err, errCreateMockData)
	_, err = db.CreateUser(context.Background(), gendb.CreateUserParams{Username: "testuser", Password: string(hashedPassword)})
	require.NoError(t, err, errCreateMockData)

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser", "testpassword", false},
		{"testuser2", "wrongpassword", true},
		{"wronguser", "testpassword3", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, err := Login(context.Background(), db, testCase.Username, testCase.Password)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Len(t, row.Token.String, tokenLen)
				u, err := db.GetUserCredentialsWithUsername(context.Background(), testCase.Username)
				require.NoError(t, err, "failed to get user by username")
				require.Equal(t, u.Token.String, row.Token.String)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	require.NoError(t, err, errCreateMockData)
	_, err = db.CreateUser(context.Background(), gendb.CreateUserParams{Username: "testuser", Password: string(hashedPassword)})
	require.NoError(t, err, errCreateMockData)

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser", "testpassword", false},
		{"wronguser", "testpassword3", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, _ := Login(context.Background(), db, testCase.Username, testCase.Password)
			err = Logout(context.Background(), db, testCase.Username)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				u, err := getUserCredentialsWithUsername(context.Background(), db, testCase.Username)
				require.NoError(t, err, "failed to get user by username")
				require.Equal(t, u.TokenExpireAt.Valid, false)
				require.Equal(t, u.Token.String, row.Token.String)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser", "testpassword", false},
		{"testuser", "testpassword", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, err := Register(context.Background(), db, testCase.Username, testCase.Password)

			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				u, err := user.GetUserByUsername(context.Background(), db, testCase.Username)
				require.NoError(t, err)

				require.Equal(t, row.ID, u.ID)
				require.Equal(t, row.Username, u.Username)
				require.Equal(t, row.Username, testCase.Username)
			}
		})
	}
}
