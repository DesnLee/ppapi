package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Tag struct {
	ID        int64              `json:"id" example:"1"`
	UserID    pgtype.UUID        `json:"user_id" swaggertype:"string" example:"0d00ec26-5f95-4f13-bdfb-a511f13aa0d2"` // UUIDv4
	Name      string             `json:"name" example:"餐饮"`
	Sign      string             `json:"sign" example:"😄"`
	Kind      string             `json:"kind" example:"expenses"`
	DeletedAt pgtype.Timestamptz `json:"deleted_at" swaggertype:"string" example:"2021-01-01T00:00:00:000+08:00"`
}

type CreateTagRequestBody struct {
	Name string `json:"name" binding:"required" example:"餐饮"`
	Sign string `json:"sign" binding:"required" example:"😄"`
	Kind string `json:"kind" binding:"required" example:"expenses"`
}
type CreateTagResponseSuccessBody = ResourceResponse[Tag]

type GetTagResponseSuccessBody = ResourceResponse[Tag]

type UpdateTagRequestBody struct {
	Name string `json:"name" example:"餐饮"`
	Sign string `json:"sign" example:"😄"`
	Kind string `json:"kind" example:"expenses"`
}
type UpdateTagResponseSuccessBody = ResourceResponse[Tag]
