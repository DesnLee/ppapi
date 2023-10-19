package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/jwt_helper"
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

	u, err := findOrCreateUser(c, body.Email)
	if err != nil {
		return
	}

	jwt, err := jwt_helper.GenerateJWT(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器错误"})
		return
	}

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
		_ = tx.Rollback()
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

// findOrCreateUser 查找或创建用户
func findOrCreateUser(c *gin.Context, email string) (sqlcExec.User, error) {
	u := sqlcExec.User{}
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("[Create Database Transaction Failed]: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器错误"})
		return u, err
	}
	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	qtx := database.Q.WithTx(tx)
	u, err = qtx.FindUserByEmail(database.DBCtx, email)

	if err != nil {
		// 如果是无记录
		if errors.Is(err, sql.ErrNoRows) {
			// 创建用户
			u, err = qtx.CreateUser(database.DBCtx, email)
			if err != nil {
				log.Println("[Create user Failed]: ", err)
				c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务错误，请稍后再试"})
				return u, err
			}
		} else {
			log.Println("[Find user by email Failed]: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务错误，请稍后再试"})
			return u, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return u, err
	}

	return u, nil
}
