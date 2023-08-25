package auth

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=15"`
	Password string `json:"password" validate:"required,min=6,max=16,lowercase,uppercase,numeric"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=5,max=15"`
	Password string `json:"password" validate:"required,min=6,max=16,lowercase,uppercase,numeric"`
}
