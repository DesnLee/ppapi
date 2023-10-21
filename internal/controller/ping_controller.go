package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/model"
)

// PingHandler godoc
//
//	@Summary		接口健康检查
//	@Description	检查接口是否正常
//	@Tags			基础
//	@Success		200
//	@Router			/ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.MsgResponse{
		Msg: "pong",
	})
}
