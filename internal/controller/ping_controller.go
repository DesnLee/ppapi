package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, struct {
		Msg string `json:"msg"`
	}{
		Msg: "pong",
	})
}
