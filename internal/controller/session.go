package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
)

type loginRequestBody struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}
type loginResponseBody struct {
	JWT string `json:"jwt"`
}

// CreateSessionHandler godoc
// @Summary      用户登录
// @Description  用户登录并获取 token
// @Accept       json
// @Produce      json
// @Param        body body loginRequestBody true "comment"
// @Success      200 {object} loginResponseBody
// @Failure      400
// @Failure      401
// @Router       /session [post]
func CreateSessionHandler(c *gin.Context) {
	body := loginRequestBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	_, err := database.Q.CheckValidationCode(database.DBCtx, sqlcExec.CheckValidationCodeParams{
		Email: body.Email,
		Code:  body.Code,
	})

	if err != nil {
		log.Println("[CheckValidationCode Failed]: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "无效的邮箱或验证码"})
		return
	}

	jwt := "JWT"
	responseBody := loginResponseBody{
		JWT: jwt,
	}

	c.JSON(http.StatusOK, responseBody)
}
