package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/controller/controller_helper"
	"ppapi.desnlee.com/internal/jwt_helper"
)

type SessionController struct{}

func (ctl *SessionController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.POST("/session", ctl.Create)
}

// Create godoc
// @Summary      用户登录
// @Description  用户登录并获取 token
// @Accept       json
// @Produce      json
// @Param        body body loginRequestBody true "comment"
// @Success      200 {object} loginResponseBody
// @Failure      400
// @Failure      401
// @Router       /session [post]
func (ctl *SessionController) Create(c *gin.Context) {
	body := struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "服务器错误"})
		return
	}

	responseBody := struct {
		JWT string `json:"jwt"`
	}{
		JWT: jwt,
	}

	c.JSON(http.StatusOK, responseBody)
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
