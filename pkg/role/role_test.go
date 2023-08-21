package role

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
)

func TestCreateRole(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err)

	testCases := []struct {
		Role        string
		ExpectedErr bool
	}{
		{"testrole1", true},
		{"nonexistingrole", false},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, err := createRole(context.Background(), db, testCase.Role)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.Role, row.Role, "expected: %s, got: %s", testCase.Role, row.Role)
			}
		})
	}
}

func TestDeleteRole(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err)

	testCases := []struct {
		Role        string
		ExpectedErr bool
	}{
		{"testrole1", false},
		{"nonexistingrole", true},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			err := deleteRole(context.Background(), db, testCase.Role)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
