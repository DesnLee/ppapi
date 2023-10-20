package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/email"
	"ppapi.desnlee.com/pkg"
)

type ValidationCodeController struct{}

func (ctl *ValidationCodeController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.POST("/validation_code", ctl.Create)
}

// Create godoc
// @Summary      邮件验证码
// @Description  发送邮件验证码
// @Accept       json
// @Produce      json
// @Param        body body getValidationCodeRequestBody true "comment"
// @Success      204
// @Failure      400
// @Failure      500
// @Router       /validation_code [post]
func (ctl *ValidationCodeController) Create(c *gin.Context) {
	body := struct {
		Email string `json:"email" binding:"required"`
	}{}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	code, err := pkg.GenerateRandomCode(6)
	if err != nil {
		log.Println("[RandCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "发送失败"})
		return
	}

	row, err := database.Q.CreateValidationCode(database.DBCtx, sqlcExec.CreateValidationCodeParams{
		Email: body.Email,
		Code:  code,
	})
	if err != nil {
		log.Println("[CreateValidationCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "发送失败"})
		return
	}

	err = email.SendValidationCode(row.Email, row.Code)
	if err != nil {
		log.Println("[SendValidationCode Failed]: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "发送失败"})
		return
	} else {
		c.Status(http.StatusNoContent)
	}
}

func (ctl *ValidationCodeController) Read(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ValidationCodeController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ValidationCodeController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ValidationCodeController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
