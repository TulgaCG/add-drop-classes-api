package user

// TODO: Add playground-go/validator tags

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	NewUsername string `json:"newUsername"`
	NewPassword string `json:"newPassword"`
}
