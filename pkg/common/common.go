package common

import "time"

const (
	ValidTime         = 1 * time.Minute
	TokenLength       = 64
	DatabaseCtxKey    = "db"
	UsernameHeaderKey = "Username"
	TokenHeaderKey    = "Token"
)

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error"`
}
