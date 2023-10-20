package jwt_helper

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}

func (m JWTClaims) Validate() error {
	if m.UserID == uuid.Nil {
		return errors.New("invalid user id")
	}
	return nil
}
