package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

const errCreateMockData = "failed to create mock data"

func TestListUsers(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	_, err = db.CreateUser(context.Background(), gendb.CreateUserParams{Username: "testuser", Password: "testpassword"})
	require.NoError(t, err, errCreateMockData)

	expectedRows, err := db.ListUsers(context.Background())
	require.NoError(t, err)

	actualRows, err := ListUsers(context.Background(), db)
	require.NoError(t, err)

	for i := range expectedRows {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require.Equal(t, expectedRows[i], actualRows[i], fmt.Sprintf("test failed on row index of %d", i))
		})
	}
}

func TestGetUser(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	_, err = db.CreateUser(context.Background(), gendb.CreateUserParams{Username: "testuser", Password: "testpassword"})
	require.NoError(t, err, errCreateMockData)

	testCases := []struct {
		Index       types.UserID
		ExpectedErr bool
	}{
		{0, true},
		{1, false},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			expected, _ := db.GetUser(context.Background(), testCase.Index)
			actual, err := GetUser(context.Background(), db, testCase.Index)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, actual)
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	_, err = db.CreateUser(context.Background(), gendb.CreateUserParams{Username: "testuser", Password: "testpassword"})
	require.NoError(t, err, errCreateMockData)

	testCases := []struct {
		Username    string
		ExpectedErr bool
	}{
		{"testuser", false},
		{"wronguser", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			expected, _ := db.GetUserByUsername(context.Background(), testCase.Username)
			actual, err := GetUserByUsername(context.Background(), db, testCase.Username)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, expected, actual)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	_, err = db.CreateUser(context.Background(), gendb.CreateUserParams{Username: "testuser", Password: "testpassword"})
	require.NoError(t, err, errCreateMockData)

	testCases := []struct {
		Index       types.UserID
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{1, "updateduser1", "updatedpassword1", false},
		{7, "wronguser", "wrongpassword", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			user, _ := db.GetUser(context.Background(), testCase.Index)
			actual, err := UpdateUser(context.Background(), db, gendb.UpdateUserParams{
				ID:       testCase.Index,
				Username: testCase.Username,
				Password: testCase.Password,
			})
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEqual(t, user, actual)

				expected, _ := db.GetUser(context.Background(), testCase.Index)
				require.Equal(t, expected.ID, actual.ID)
				require.Equal(t, expected.Username, actual.Username)
			}
		})
	}
}
