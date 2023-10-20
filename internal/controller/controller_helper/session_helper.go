package controller_helper

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/model"
)

// CheckAndUseValidationCode 校验并使用验证码
func CheckAndUseValidationCode(c *gin.Context, email, code string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("[Create Database Transaction Failed]: ", err)
		c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "服务器错误"})
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
		c.JSON(http.StatusUnauthorized, model.MsgResponse{Msg: "无效的邮箱或验证码"})
		return err
	}

	// 使用验证码
	if _, err = qtx.UseValidationCode(database.DBCtx, r.ID); err != nil {
		log.Println("[UseValidationCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "服务器错误"})
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// FindOrCreateUserByEmail 查找或创建用户
func FindOrCreateUserByEmail(c *gin.Context, email string) (sqlcExec.User, error) {
	u := sqlcExec.User{}
	tx, err := database.DB.Begin()
	if err != nil {
		log.Println("[Create Database Transaction Failed]: ", err)
		c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "服务器错误"})
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
				c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "服务器错误"})
				return u, err
			}
		} else {
			log.Println("[Find user by email Failed]: ", err)
			c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "服务器错误"})
			return u, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return u, err
	}

	return u, nil
}
