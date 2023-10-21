package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/email"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

type ValidationCodeController struct{}

func (ctl *ValidationCodeController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.POST("/validation_code", ctl.Create)
}

// Create godoc
//
//	@Summary		邮件验证码
//	@Description	发送邮件验证码
//	@Tags			登录
//	@Accept			json
//	@Produce		json
//	@Param			body	body	model.ValidationCodeRequestBody	true	"传入接收验证码的邮箱，未注册将会自动注册"
//	@Success		204
//	@Failure		400	{object}	model.MsgResponse	"参数错误"
//	@Failure		500	{object}	model.MsgResponse	"服务器错误"
//	@Router			/validation_code [post]
func (ctl *ValidationCodeController) Create(c *gin.Context) {
	body := model.ValidationCodeRequestBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{Msg: "参数错误"})
		return
	}

	code, err := pkg.GenerateRandomCode(6)
	if err != nil {
		log.Println("[RandCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "发送失败"})
		return
	}

	row, err := database.Q.CreateValidationCode(database.DBCtx, sqlcExec.CreateValidationCodeParams{
		Email: body.Email,
		Code:  code,
	})
	if err != nil {
		log.Println("[CreateValidationCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "发送失败"})
		return
	}

	err = email.SendValidationCode(row.Email, row.Code)
	if err != nil {
		log.Println("[SendValidationCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "发送失败"})
		return
	} else {
		c.Status(http.StatusNoContent)
	}
}

func (ctl *ValidationCodeController) Read(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ValidationCodeController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ValidationCodeController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ValidationCodeController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
