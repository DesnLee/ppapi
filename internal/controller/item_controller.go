package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/internal/controller/controller_helper"
	"ppapi.desnlee.com/internal/middleware"
	"ppapi.desnlee.com/internal/model"
)

type ItemController struct{}

func (ctl *ItemController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.Use(middleware.JWTMiddleware())
	v1.POST("/item", ctl.Create)
}

// Create godoc
//
//	@Summary		新建记账条目
//	@Description	新建记账条目
//	@Tags			账单
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.CreateItemRequestBody			true	"传入记账信息"
//	@Success		200		{object}	model.CreateItemResponseSuccessBody	"成功创建记账条目"
//	@Failure		401		{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		422		{object}	model.MsgResponse					"参数错误"
//	@Failure		500		{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/item [post]
func (ctl *ItemController) Create(c *gin.Context) {
	userID := c.MustGet("userID").(pgtype.UUID)
	body := model.CreateItemRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{
			Msg: "参数错误",
		})
		return
	}

	if err := controller_helper.ValidateCreateItemRequestBody(userID, &body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.MsgResponse{
			Msg: err.Error(),
		})
		return
	}

	r, err := controller_helper.CreateItem(userID, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.MsgResponse{
			Msg: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.CreateItemResponseSuccessBody{Resource: r})
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
