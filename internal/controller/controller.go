package controller

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Register(r *gin.RouterGroup)

	Create(c *gin.Context)
	Read(c *gin.Context)
	ReadMulti(c *gin.Context)
	Update(c *gin.Context)
	Destroy(c *gin.Context)
}
