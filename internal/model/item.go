package model

import (
	"time"
)

type CreateItemRequestBody struct {
	Amount     int64     `json:"amount" binding:"required"`
	Kind       string    `json:"kind" binding:"required"`
	HappenedAt time.Time `json:"happened_at" binding:"required"`
	TagIDs     []int64   `json:"tag_ids" binding:"required"`
}

type CreateItemResponseBody struct {
	ID         int64     `json:"id"`
	Amount     int64     `json:"amount"`
	Kind       string    `json:"kind"`
	HappenedAt time.Time `json:"happened_at"`
}
