package model

import (
	"ppapi.desnlee.com/db/sqlcExec"
)

type TagResponse struct {
	ID   int64         `json:"id"`
	Name string        `json:"name"`
	Sign string        `json:"sign"`
	Kind sqlcExec.Kind `json:"kind"`
}
