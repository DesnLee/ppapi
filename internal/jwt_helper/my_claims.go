package jwt_helper

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type JWTClaims struct {
	UserID pgtype.UUID `json:"userId"`
	jwt.RegisteredClaims
}

func (m JWTClaims) Validate() error {
	if !m.UserID.Valid {
		return errors.New("invalid user id")
	}
	return nil
}
