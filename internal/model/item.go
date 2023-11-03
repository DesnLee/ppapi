package model

import (
	"ppapi.desnlee.com/db/sqlcExec"
)

type CreateItemRequestBody struct {
	Amount     int64         `json:"amount" binding:"required" example:"100"` // 金额单位为（分）
	Kind       sqlcExec.Kind `json:"kind" binding:"required" swaggertype:"string" enums:"income,expenses" example:"expenses"`
	HappenedAt string        `json:"happened_at" binding:"required" example:"2023-01-01"`
	TagIDs     []int64       `json:"tag_ids" binding:"required" example:"1"`
}

type CreateItemResponseData struct {
	ID         int64         `json:"id" example:"1"`
	Amount     int64         `json:"amount" example:"100"`
	Kind       sqlcExec.Kind `json:"kind" swaggertype:"string" example:"expenses" enums:"income,expenses"`
	HappenedAt string        `json:"happened_at"  example:"2023-01-01"`
	TagIDs     []int64       `json:"tag_ids" example:"1"`
}
type CreateItemResponseSuccessBody = ResourceResponse[CreateItemResponseData]

type GetItemsRequestBody struct {
	Page           int64  `form:"page" example:"1"`
	Size           int64  `form:"size" example:"10"`
	HappenedBefore string `form:"happened_before" example:"2023-01-01"`
	HappenedAfter  string `form:"happened_after" example:"2023-01-01"`
}
type GetItemsResponseData struct {
	CreateItemResponseData
	Tags []TagResponse `json:"tags"`
}
type GetItemsResponseSuccessBody = ResourcesResponse[GetItemsResponseData]
