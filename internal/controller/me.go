package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/jwt_helper"
	"ppapi.desnlee.com/internal/model"
)

type MeController struct{}

func (ctl *MeController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.GET("/me", ctl.Read)
}

func (ctl *MeController) Create(c *gin.Context) {
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
//	@Router			/api/v1/me [get]
func (ctl *MeController) Read(c *gin.Context) {
	authStr := c.GetHeader("Authorization")
	if authStr == "" {
		c.JSON(401, model.MsgResponse{
			Msg: "未携带 token",
		})
		return
	}

	parts := strings.SplitN(authStr, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(401, model.MsgResponse{
			Msg: "token 格式错误",
		})
		return
	}

	claims, err := jwt_helper.ParseJWT(parts[1])
	if err != nil {
		c.JSON(401, model.MsgResponse{
			Msg: "token 无效",
		})
		return
	}

	userID := claims.UserID
	if userID == uuid.Nil {
		c.JSON(401, model.MsgResponse{
			Msg: "token 无效",
		})
		return
	}

	u, err := database.Q.FindUserByID(database.DBCtx, userID)
	if err != nil {
		log.Println("ERR: [Find User By ID Failed]: ", err)
		c.JSON(500, model.MsgResponse{
			Msg: "服务器错误",
		})
		return
	}

	c.JSON(http.StatusOK, model.ResourceResponse[model.MeResponseBody]{
		Resource: model.MeResponseBody{
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
