package token

import (
	"time"

	"github.com/google/uuid"
)

type ITokenMaker interface {
	// CreateToken creates a new token for a specific email and duration
	CreateToken(uuid.UUID, time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(string) (*Payload, error)
}
