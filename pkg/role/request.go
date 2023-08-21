package role

import "github.com/TulgaCG/add-drop-classes-api/pkg/types"

type AddRoleRequest struct {
	UserID types.UserID `json:"userid"`
	RoleID types.RoleID `json:"roleid"`
}
