package model

import (
	"ppapi.desnlee.com/db/sqlcExec"
)

type TagResponse struct {
	ID   int64         `json:"id" example:"1"`
	Name string        `json:"name" example:"餐饮"`
	Sign string        `json:"sign" example:"😄"`
	Kind sqlcExec.Kind `json:"kind" swaggertype:"string" example:"expenses" enums:"expenses,incomes"`
}
