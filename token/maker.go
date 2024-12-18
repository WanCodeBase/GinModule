package token

import "time"

type Maker interface {
	// CreateToken returns token
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken verify token is valid
	VerifyToken(token string) (*Payload, error)
}
