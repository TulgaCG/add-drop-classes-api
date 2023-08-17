package key

import "time"

const (
	ValidTime         = 1 * time.Minute
	TokenLength       = 64
	DatabaseCtxKey    = "db"
	UsernameHeaderKey = "Username"
	TokenHeaderKey    = "Token"
)
