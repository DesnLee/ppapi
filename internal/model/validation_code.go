package model

type ValidationCodeRequestBody struct {
	Email string `json:"email" binding:"required"`
}
