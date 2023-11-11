package model

import (
	"ppapi.desnlee.com/db/sqlcExec"
)

type Tag struct {
	ID     int64         `json:"id" example:"1"`
	UserID int64         `json:"user_id" example:"1"`
	Name   string        `json:"name" example:"餐饮"`
	Sign   string        `json:"sign" example:"😄"`
	Kind   sqlcExec.Kind `json:"kind" swaggertype:"string" example:"expenses" enums:"income,expenses"`
	// DeletedAt string        `json:"deleted_at" example:"2021-01-01T00:00:00+08:00"`
}

type CreateTagRequestBody struct {
	Name string        `json:"name" example:"餐饮"`
	Sign string        `json:"sign" example:"😄"`
	Kind sqlcExec.Kind `json:"kind" swaggertype:"string" example:"expenses" enums:"income,expenses"`
}
type CreateTagResponseSuccessBody = ResourceResponse[Tag]
