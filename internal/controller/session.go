package controller

import (
	"database/sql"
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

	// 使用数据库事务校验并使用验证码
	if err := checkAndUseValidationCode(c, body.Email, body.Code); err != nil {
		return
	}

	jwt := "JWT"
	responseBody := loginResponseBody{
		JWT: jwt,
	}

	c.JSON(http.StatusOK, responseBody)
}

// checkAndUseValidationCode 校验并使用验证码
func checkAndUseValidationCode(c *gin.Context, email, code string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("[Create Database Transaction Failed]: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器错误"})
		return err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Println("[Rollback Transaction Failed]: ", err)
		}
	}(tx)

	qtx := database.Q.WithTx(tx)
	r, err := qtx.CheckValidationCode(database.DBCtx, sqlcExec.CheckValidationCodeParams{
		Email: email,
		Code:  code,
	})

	if err != nil {
		log.Println("[CheckValidationCode Failed]: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "无效的邮箱或验证码"})
		return err
	}

	// 使用验证码
	if _, err = qtx.UseValidationCode(database.DBCtx, r.ID); err != nil {
		log.Println("[UseValidationCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器错误"})
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
