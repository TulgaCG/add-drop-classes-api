package lecture

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func TestGetLecturesFromUser(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	testCases := []struct {
		UserID      types.UserID
		ExpectedErr bool
	}{
		{types.UserID(1), false},
		{types.UserID(4), true},
		{types.UserID(10), true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			actual, err := GetLecturesFromUser(context.Background(), db, testCase.UserID)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				expected, err := db.GetLecture(context.Background(), types.LectureID(testCase.UserID))
				require.NoError(t, err)
				require.Equal(t, expected.Code, actual[0].Code)
				require.Equal(t, expected.Name, actual[0].Name)
				require.Equal(t, expected.Credit, actual[0].Credit)
			}
		})
	}
}

func TestAddLectureToUser(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	testCases := []struct {
		Username    string
		LectureCode string
		ExpectedErr bool
	}{
		{"testuser1", "ADA102", false},
		{"testuser2", "ADA102", true},
		{"testuser10", "ADA102", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			_, err := AddLectureToUser(context.Background(), db, testCase.Username, testCase.LectureCode)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRemoveLectureFromUser(t *testing.T) {
	db, closeFn, err := database.NewTestDB(context.Background())
	require.NoError(t, err)
	defer closeFn(t)

	testCases := []struct {
		UserID      types.UserID
		LectureID   types.LectureID
		ExpectedErr bool
	}{
		{1, 1, false},
		{1, 2, true},
		{10, 1, true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			err := RemoveLectureFromUser(context.Background(), db, testCase.UserID, testCase.LectureID)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
