package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/middleware"
	"ppapi.desnlee.com/internal/model"
)

type ItemController struct{}

func (ctl *ItemController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.Use(middleware.JWTMiddleware())
	v1.POST("/item", ctl.Create)
}

type createItemResponseSuccessData struct {
	ID         int64     `json:"id"`
	Amount     int64     `json:"amount"`
	Kind       string    `json:"kind"`
	HappenedAt time.Time `json:"happened_at"`
	TagIDs     []int64   `json:"tag_ids"`
}
type createItemResponseSuccessBodyForDoc = model.ResourceResponse[createItemResponseSuccessData]
type createItemRequestBodyForDoc struct {
	Amount     int64     `json:"amount" binding:"required"`
	Kind       string    `json:"kind" binding:"required"`
	HappenedAt time.Time `json:"happened_at" binding:"required"`
	TagIDs     []int64   `json:"tag_ids" binding:"required"`
}

// Create godoc
//
//	@Summary		新建记账条目
//	@Description	新建记账条目
//	@Tags			账单
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"token字符串，格式 `Bearer {token}`"
//	@Param			body			body		createItemRequestBodyForDoc			true	"传入记账信息"
//	@Success		200				{object}	createItemResponseSuccessBodyForDoc	"成功创建记账条目"
//	@Failure		401				{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		500				{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/item [post]
func (ctl *ItemController) Create(c *gin.Context) {
	userID := c.MustGet("userID").(pgtype.UUID)
	body := model.CreateItemRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{
			Msg: "参数错误",
		})
		return
	}

	tags, err := database.Q.FindTagsByIDs(database.DBCtx, body.TagIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.MsgResponse{
			Msg: "服务器错误",
		})
		return
	}
	fmt.Println(tags)
	fmt.Println(userID)
}

func (ctl *ItemController) Read(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ItemController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ItemController) Update(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (ctl *ItemController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
