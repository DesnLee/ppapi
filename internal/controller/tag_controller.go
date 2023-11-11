package controller

import (
	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/internal/middleware"
)

type TagController struct{}

func (ctl *TagController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.Use(middleware.JWTMiddleware())
	v1.POST("/tags", ctl.Create)
}

func (ctl *TagController) Create(c *gin.Context) {
	// TODO implement me
	panic("implement me")
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
