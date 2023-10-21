package model

import (
	"github.com/google/uuid"
)

type MeResponseBody struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}
