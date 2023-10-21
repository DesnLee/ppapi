package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/controller/controller_helper"
	"ppapi.desnlee.com/internal/jwt_helper"
	"ppapi.desnlee.com/internal/model"
)

type SessionController struct{}

func (ctl *SessionController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.POST("/session", ctl.Create)
}

// Create godoc
//
//	@Summary		用户登录
//	@Description	用户登录并获取 token
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.SessionRequestBody						true	"传入邮箱和验证码"
//	@Success		200		{object}	model.DataResponse[model.SessionResponseBody]	"成功获取到 jwt"
//	@Failure		400		{object}	model.MsgResponse								"参数错误"
//	@Failure		401		{object}	model.MsgResponse								"验证码错误"
//	@Failure		500		{object}	model.MsgResponse								"服务器错误"
//	@Router			/session [post]
func (ctl *SessionController) Create(c *gin.Context) {
	body := model.SessionRequestBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{Msg: "参数错误"})
		return
	}

	// 使用数据库事务校验并使用验证码
	if err := controller_helper.CheckAndUseValidationCode(c, body.Email, body.Code); err != nil {
		return
	}

	u, err := controller_helper.FindOrCreateUserByEmail(c, body.Email)
	if err != nil {
		return
	}

	jwt, err := jwt_helper.GenerateJWT(u.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.MsgResponse{Msg: "参数错误"})
		return
	}

	c.JSON(http.StatusOK, model.SessionResponseBody{
		JWT: jwt,
	})
}

func (ctl *SessionController) Read(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *SessionController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *SessionController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *SessionController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
