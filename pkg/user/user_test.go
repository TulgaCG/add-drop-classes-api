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

func TestCreateUser(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

	testCases := []struct {
		Username    string
		Password    string
		ExpectedErr bool
	}{
		{"testuser1", "testpassword1", true},
		{"testuser10", "testpassword10", false},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, err := createUser(context.Background(), db, testCase.Username, testCase.Password)

			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				user, err := getUserByUsername(context.Background(), db, testCase.Username)
				require.NoError(t, err)

				require.Equal(t, row.ID, user.ID)
				require.Equal(t, row.Username, user.Username)
				require.Equal(t, row.Username, testCase.Username)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

	expectedRows, err := db.ListUsers(context.Background())
	require.NoError(t, err)

	actualRows, err := listUsers(context.Background(), db)
	require.NoError(t, err)

	for i := range expectedRows {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			require.Equal(t, expectedRows[i], actualRows[i], fmt.Sprintf("test failed on row index of %d", i))
		})
	}
}

func TestGetUser(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

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
			actual, err := getUser(context.Background(), db, testCase.Index)
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
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

	testCases := []struct {
		Username    string
		ExpectedErr bool
	}{
		{"testuser1", false},
		{"wronguser", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			expected, _ := db.GetUserByUsername(context.Background(), testCase.Username)
			actual, err := getUserByUsername(context.Background(), db, testCase.Username)
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
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err, "failed to create test db")

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
			actual, err := updateUser(context.Background(), db, gendb.UpdateUserParams{
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
