package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/model"
)

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, model.MsgResponse{
		Msg: "pong",
	})
}
