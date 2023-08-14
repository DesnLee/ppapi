package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type requestBody struct {
	Email string `json:"email" binding:"required"`
}

// SendValidationCodeHandler godoc
// @Summary      邮件验证码
// @Description  发送邮件验证码
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /validation_code [post]
func SendValidationCodeHandler(c *gin.Context) {
	body := requestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}

	fmt.Println(body.Email)
}
