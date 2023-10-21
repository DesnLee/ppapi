package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MeResponseBody struct {
	ID    pgtype.UUID `json:"id"`
	Email string      `json:"email"`
	Name  string      `json:"name"`
}
