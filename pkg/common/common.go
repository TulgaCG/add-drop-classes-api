package common

import (
	"github.com/TulgaCG/add-drop-classes-api/pkg/types"
)

const (
	DatabaseCtxKey    = "db"
	LogCtxKey         = "log"
	RolesCtxKey       = "roles"
	UsernameHeaderKey = "Username"
	DefaultRole       = types.RoleID(3)
)
