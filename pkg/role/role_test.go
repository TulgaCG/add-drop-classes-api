package role

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TulgaCG/add-drop-classes-api/pkg/database"
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

func TestAddRoleToUser(t *testing.T) {
	db, err := database.NewTestDb(context.Background())
	require.NoError(t, err)

	testCases := []struct {
		UserID      types.UserID
		RoleID      types.RoleID
		ExpectedErr bool
	}{
		{1, 1, false},
		{1, 1, true},
		{1, 4, false},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			row, err := addRoleToUser(context.Background(), db, testCase.UserID, testCase.RoleID)
			if testCase.ExpectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.UserID, row.UserID)
			}
		})
	}
}
