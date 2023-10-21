package controller

import (
	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/model"
)

type Me struct{}

func (ctl *Me) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.POST("/me", ctl.Create)
}

func (ctl *Me) Create(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

type getResponseBody = model.ResourceResponse[model.MeResponseBody]

// Read godoc
//
//	@Summary		获取当前用户
//	@Description	获取当前用户的基本信息
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"token字符串，格式 `Bearer {token}`"
//	@Success		200				{object}	getResponseBody		"成功获取到用户信息"
//	@Failure		401				{object}	model.MsgResponse	"未授权，token 无效"
//	@Failure		500				{object}	model.MsgResponse	"服务器错误"
//	@Router			/v1/me [get]
func (ctl *Me) Read(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *Me) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *Me) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *Me) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
