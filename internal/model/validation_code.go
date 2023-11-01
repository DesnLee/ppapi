package model

type ValidationCodeRequestBody struct {
	Email string `json:"email" binding:"required" example:"test@qq.com"` // 接收验证码的邮箱
}
