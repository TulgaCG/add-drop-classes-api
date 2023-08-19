package auth

// TODO: Add playground-go/validator tags

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
