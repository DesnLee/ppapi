package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/controller/controller_helper"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/middleware"
	"ppapi.desnlee.com/internal/model"
)

type TagController struct{}

func (ctl *TagController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.Use(middleware.JWTMiddleware())
	v1.POST("/tags", ctl.Create)
}

// Create godoc
//
//	@Summary		新建标签
//	@Description	新建标签
//	@Tags			标签
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.CreateTagRequestBody			true	"传入标签信息"
//	@Success		200		{object}	model.CreateTagResponseSuccessBody	"成功创建标签"
//	@Failure		401		{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		422		{object}	model.MsgResponse					"参数错误"
//	@Failure		500		{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/tags [post]
func (ctl *TagController) Create(c *gin.Context) {
	userID := c.MustGet("userID").(pgtype.UUID)
	body := model.CreateTagRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{
			Msg: "参数错误",
		})
		return
	}

	if err := controller_helper.ValidateCreateTagRequestBody(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.MsgResponse{
			Msg: err.Error(),
		})
		return
	}

	r, err := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: userID,
		Name:   body.Name,
		Sign:   body.Sign,
		Kind:   body.Kind,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.MsgResponse{
			Msg: "服务器错误，创建标签失败",
		})
		return
	}

	c.JSON(http.StatusOK, model.CreateTagResponseSuccessBody{Resource: model.Tag{
		ID:     r.ID,
		UserID: r.UserID,
		Name:   r.Name,
		Sign:   r.Sign,
		Kind:   r.Kind,
	}})
}

func (ctl *TagController) Read(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *TagController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *TagController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *TagController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
