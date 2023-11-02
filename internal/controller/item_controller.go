package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgtype"
    "ppapi.desnlee.com/internal/controller/controller_helper"
    "ppapi.desnlee.com/internal/middleware"
    "ppapi.desnlee.com/internal/model"
)

type ItemController struct{}

func (ctl *ItemController) Register(g *gin.RouterGroup) {
    v1 := g.Group("/v1")
    v1.Use(middleware.JWTMiddleware())
    v1.POST("/items", ctl.Create)
    v1.GET("/items", ctl.ReadMulti)
}

// Create godoc
//
//	@Summary		新建记账条目
//	@Description	新建记账条目
//	@Tags			账单
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.CreateItemRequestBody			true	"传入记账信息"
//	@Success		200		{object}	model.CreateItemResponseSuccessBody	"成功创建记账条目"
//	@Failure		401		{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		422		{object}	model.MsgResponse					"参数错误"
//	@Failure		500		{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/items [post]
func (ctl *ItemController) Create(c *gin.Context) {
    userID := c.MustGet("userID").(pgtype.UUID)
    body := model.CreateItemRequestBody{}
    if err := c.ShouldBindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, model.MsgResponse{
            Msg: "参数错误",
        })
        return
    }

    if err := controller_helper.ValidateCreateItemRequestBody(userID, &body); err != nil {
        c.JSON(http.StatusUnprocessableEntity, model.MsgResponse{
            Msg: err.Error(),
        })
        return
    }

    r, err := controller_helper.CreateItem(userID, &body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.MsgResponse{
            Msg: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, model.CreateItemResponseSuccessBody{Resource: r})
}

func (ctl *ItemController) Read(c *gin.Context) {
    // TODO implement me
    panic("implement me")
}

// ReadMulti godoc
//
//	@Summary		查询记账条目
//	@Description	查询多条记账条目
//	@Tags			账单
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			query	query		model.GetItemsRequestBody			true	"传入查询条件"
//	@Success		200		{object}	model.GetItemsResponseSuccessBody	"成功查询记账条目"
//	@Failure		401		{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		422		{object}	model.MsgResponse					"参数错误"
//	@Failure		500		{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/items [get]
func (ctl *ItemController) ReadMulti(c *gin.Context) {
    userID := c.MustGet("userID").(pgtype.UUID)
    query := model.GetItemsRequestBody{}
    if err := c.ShouldBindQuery(&query); err != nil {
        c.JSON(http.StatusBadRequest, model.MsgResponse{
            Msg: "参数错误",
        })
        return
    }

    res, err := controller_helper.GetAndCountItemsByUserID(userID, query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, model.MsgResponse{
            Msg: err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, res)
}

func (ctl *ItemController) Update(c *gin.Context) {
    // TODO implement me
    panic("implement me")
}

func (ctl *ItemController) Destroy(c *gin.Context) {
    // TODO implement me
    panic("implement me")
}
