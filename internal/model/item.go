package model

import (
	"time"
)

type CreateItemRequestBody struct {
	Amount     int       `json:"amount" binding:"required"`
	Kind       string    `json:"kind" binding:"required"`
	HappenedAt time.Time `json:"happened_at" binding:"required"`
	TagIDs     []int     `json:"tag_ids" binding:"required"`
}

type CreateItemResponseBody struct {
	ID         int       `json:"id"`
	Amount     int       `json:"amount"`
	Kind       string    `json:"kind"`
	HappenedAt time.Time `json:"happened_at"`
	TagIDs     []int     `json:"tag_ids"`
}
