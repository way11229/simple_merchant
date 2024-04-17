package auth_token_maker

import "time"

type AuthTokenMaker interface {
	CreateToken(uniqueCode string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
