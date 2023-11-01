package model

type SessionRequestBody struct {
	Email string `json:"email" binding:"required" example:"test@qq.com"` // 登录邮箱
	Code  string `json:"code" binding:"required" example:"123456"`       // 验证码
}

type SessionResponseBody struct {
	JWT string `json:"jwt" example:"jIxCxLH-17fUbMgl5G4bEK0dJ-JeEbj36n8pLWT7GGCsos0wfXFs6Xk8jM2R7g7zifsVjYZZuZRizU_DtIrd52caNRQqFgaNetWLTTArlNNdMsnCwndWshEGh7JC0e74lnrZs5wX75UjXscqZfmulXEBpYcr8MVeYKOhgctfPIBc-dH9Qs4KgUqS55S3fY3Go_OmssOnw5ErdBt4dr_NLfwXrw=="` // jwt 经过 sha 加密后的字符串
}
