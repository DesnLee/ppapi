package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/middleware"
	"ppapi.desnlee.com/internal/model"
)

type MeController struct{}

func (ctl *MeController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.Use(middleware.JWTMiddleware())
	v1.GET("/me", ctl.Read)
}

func (ctl *MeController) Create(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

// Read godoc
//
//	@Summary		获取当前用户
//	@Description	获取当前用户的基本信息
//	@Tags			用户
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	model.MeResponseSuccessBody	"成功获取到用户信息"
//	@Failure		401	{object}	model.MsgResponse			"未授权，token 无效"
//	@Failure		500	{object}	model.MsgResponse			"服务器错误"
//	@Router			/api/v1/me [get]
func (ctl *MeController) Read(c *gin.Context) {
	userID := c.MustGet("userID").(pgtype.UUID)
	u, err := database.Q.FindUserByID(database.DBCtx, userID)
	if err != nil {
		log.Println("ERR: [Find User By ID Failed]: ", err)
		c.JSON(http.StatusInternalServerError, model.MsgResponse{
			Msg: "服务器错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.ResourceResponse[model.MeResponseData]{
		Resource: model.MeResponseData{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		},
	})
}

func (ctl *MeController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *MeController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *MeController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
