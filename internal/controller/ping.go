package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg string `json:"msg"`
}

// PingHandler godoc
// @Summary      测试接口
// @Description  测试 API 服务是否正常
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500
// @Router       /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Msg: "pong",
	})
}
