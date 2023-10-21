package controller

import (
	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/middleware"
	"ppapi.desnlee.com/internal/model"
)

type ItemController struct{}

func (ctl *ItemController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.Use(middleware.JWTMiddleware())
	v1.POST("/item", ctl.Create)
}

type createItemResponseSuccessBody = model.ResourceResponse[model.CreateItemResponseBody]

// Create godoc
//
//	@Summary		新建记账条目
//	@Description	新建记账条目
//	@Tags			账单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"token字符串，格式 `Bearer {token}`"
//	@Param			body			body		model.CreateItemRequestBody		true	"传入条目信息"
//	@Success		200				{object}	createItemResponseSuccessBody	"成功获取到信息"
//	@Failure		401				{object}	model.MsgResponse				"未授权，token 无效"
//	@Failure		500				{object}	model.MsgResponse				"服务器错误"
//	@Router			/api/v1/item [post]
func (ctl *ItemController) Create(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ItemController) Read(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ItemController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ItemController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ItemController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
