package model

import (
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/db/sqlcExec"
)

type Tag struct {
	ID     int64         `json:"id" example:"1"`
	UserID pgtype.UUID   `json:"user_id" example:"0d00ec26-5f95-4f13-bdfb-a511f13aa0d2"` // UUIDv4
	Name   string        `json:"name" example:"餐饮"`
	Sign   string        `json:"sign" example:"😄"`
	Kind   sqlcExec.Kind `json:"kind" swaggertype:"string" example:"expenses" enums:"income,expenses"`
	// DeletedAt string        `json:"deleted_at" example:"2021-01-01T00:00:00+08:00"`
}

type CreateTagRequestBody struct {
	Name string        `json:"name" binding:"required" example:"餐饮"`
	Sign string        `json:"sign" binding:"required" example:"😄"`
	Kind sqlcExec.Kind `json:"kind" binding:"required" swaggertype:"string" example:"expenses" enums:"income,expenses"`
}
type CreateTagResponseSuccessBody = ResourceResponse[Tag]
