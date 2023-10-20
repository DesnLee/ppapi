package model

type SessionRequestBody struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type SessionResponseBody struct {
	JWT string `json:"jwt"`
}
