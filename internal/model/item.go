package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateItemRequestBody struct {
	Amount     int64              `json:"amount" binding:"required"`
	Kind       string             `json:"kind" binding:"required"`
	HappenedAt pgtype.Timestamptz `json:"happened_at" binding:"required"`
	TagIDs     []int64            `json:"tag_ids" binding:"required"`
}

type CreateItemResponseBody struct {
	ID         int64              `json:"id"`
	Amount     int64              `json:"amount"`
	Kind       string             `json:"kind"`
	HappenedAt pgtype.Timestamptz `json:"happened_at"`
	TagIDs     []int64            `json:"tag_ids"`
}
