package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TulgaCG/add-drop-classes-api/pkg/gendb"
)

func TestNew(t *testing.T) {
	db, err := New(context.Background(), ":memory:")
	require.NoError(t, err)

	expected := gendb.CreateUserParams{
		Username: "testusernamex",
		Password: "testpasswordx",
	}
	_, err = db.CreateUser(context.Background(), expected)
	require.NoError(t, err)

	actual, err := db.GetUserCredentialsWithUsername(context.Background(), expected.Username)
	require.NoError(t, err)

	require.Equal(t, expected.Username, actual.Username)
	require.Equal(t, expected.Password, actual.Password)
	require.Equal(t, 1, int(actual.ID))
}
