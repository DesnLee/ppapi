package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MeResponseData struct {
	ID    pgtype.UUID `json:"id" example:"0d00ec26-5f95-4f13-bdfb-a511f13aa0d2"` // UUIDv4
	Email string      `json:"email" example:"test@qq.com"`
	Name  string      `json:"name" example:"test user"`
}
type MeResponseSuccessBody = ResourceResponse[MeResponseData]
