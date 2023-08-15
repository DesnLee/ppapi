package controller

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/email"
)

type requestBody struct {
	Email string `json:"email" binding:"required"`
}

// SendValidationCodeHandler godoc
// @Summary      邮件验证码
// @Description  发送邮件验证码
// @Accept       json
// @Produce      json
// @Param        body body requestBody true "comment"
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

	if err := email.SendValidationCode(body.Email, "123456"); err != nil {
		log.Println("[SendValidationCode Failed]: ", err)
		c.JSON(500, gin.H{"msg": "发送失败"})
		return
	}
}
