package user

type UpdateUserRequest struct {
	Username    string `json:"username" validate:"required,min=5,max=15"`
	NewUsername string `json:"newUsername" validate:"required,min=5,max=15"`
	NewPassword string `json:"newPassword" validate:"required,lowercase,uppercase,numeric,min=6,max=16"`
}
