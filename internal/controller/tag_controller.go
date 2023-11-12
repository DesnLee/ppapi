package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/controller/controller_helper"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/middleware"
	"ppapi.desnlee.com/internal/model"
)

type TagController struct{}

func (ctl *TagController) Register(g *gin.RouterGroup) {
	v1 := g.Group("/v1")
	v1.Use(middleware.JWTMiddleware())
	v1.POST("/tags", ctl.Create)
	v1.GET("/tags/:id", ctl.Read)
	v1.PATCH("/tags/:id", ctl.Update)
}

// Create godoc
//
//	@Summary		新建标签
//	@Description	新建标签
//	@Tags			标签
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		model.CreateTagRequestBody			true	"传入标签信息"
//	@Success		200		{object}	model.CreateTagResponseSuccessBody	"成功创建标签"
//	@Failure		401		{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		422		{object}	model.MsgResponse					"参数错误"
//	@Failure		500		{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/tags [post]
func (ctl *TagController) Create(c *gin.Context) {
	userID := c.MustGet("userID").(pgtype.UUID)
	body := model.CreateTagRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{
			Msg: "参数错误",
		})
		return
	}

	if err := controller_helper.ValidateCreateTagRequestBody(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.MsgResponse{
			Msg: err.Error(),
		})
		return
	}

	r, err := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: userID,
		Name:   body.Name,
		Sign:   body.Sign,
		Kind:   body.Kind,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.MsgResponse{
			Msg: "服务器错误，创建标签失败",
		})
		return
	}

	c.JSON(http.StatusOK, model.CreateTagResponseSuccessBody{Resource: model.Tag{
		ID:        r.ID,
		UserID:    r.UserID,
		Name:      r.Name,
		Sign:      r.Sign,
		Kind:      r.Kind,
		DeletedAt: r.DeletedAt,
	}})
}

// Read godoc
//
//	@Summary		查询标签
//	@Description	查询标签
//	@Tags			标签
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int									true	"标签 ID"
//	@Success		200	{object}	model.CreateTagResponseSuccessBody	"成功查询到标签"
//	@Failure		401	{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		422	{object}	model.MsgResponse					"参数错误"
//	@Failure		500	{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/tags/{id} [get]
func (ctl *TagController) Read(c *gin.Context) {
	userID := c.MustGet("userID").(pgtype.UUID)
	id, ok := strconv.Atoi(c.Param("id"))
	if ok != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{
			Msg: "id 参数错误",
		})
		return
	}

	r, err := database.Q.FindTagByID(database.DBCtx, sqlcExec.FindTagByIDParams{
		UserID: userID,
		ID:     int64(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, model.MsgResponse{
				Msg: "标签不存在",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, model.MsgResponse{
				Msg: "服务器错误",
			})
			return
		}
	}

	c.JSON(http.StatusOK, model.GetTagResponseSuccessBody{Resource: model.Tag{
		ID:        r.ID,
		UserID:    r.UserID,
		Name:      r.Name,
		Sign:      r.Sign,
		Kind:      r.Kind,
		DeletedAt: r.DeletedAt,
	}})
}

func (ctl *TagController) ReadMulti(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}

// Update godoc
//
//	@Summary		更新标签
//	@Description	更新标签
//	@Tags			标签
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int									true	"标签 ID"
//	@Param			body	body		model.UpdateTagRequestBody			true	"传入标签信息"
//	@Success		200		{object}	model.UpdateTagResponseSuccessBody	"成功更新标签"
//	@Failure		401		{object}	model.MsgResponse					"未授权，token 无效"
//	@Failure		422		{object}	model.MsgResponse					"参数错误"
//	@Failure		500		{object}	model.MsgResponse					"服务器错误"
//	@Router			/api/v1/tags/{id} [patch]
func (ctl *TagController) Update(c *gin.Context) {
	userID := c.MustGet("userID").(pgtype.UUID)

	id, ok := strconv.Atoi(c.Param("id"))
	if ok != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{
			Msg: "id 参数错误",
		})
		return
	}

	body := model.UpdateTagRequestBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, model.MsgResponse{
			Msg: "参数错误",
		})
		return
	}
	if err := controller_helper.ValidateUpdateTagRequestBody(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, model.MsgResponse{
			Msg: err.Error(),
		})
		return
	}

	queryParams := sqlcExec.UpdateTagByIDParams{
		UserID: userID,
		ID:     int64(id),
		Name:   body.Name,
		Sign:   body.Sign,
		Kind:   body.Kind,
	}
	// if body.Name.Valid {
	//     queryParams.Name = body.Name.String
	// }
	// if body.Sign.Valid {
	//     queryParams.Sign = body.Sign.String
	// }
	// if body.Kind.Valid {
	//     queryParams.Kind = body.Kind.String
	// }
	r, err := database.Q.UpdateTagByID(database.DBCtx, queryParams)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, model.MsgResponse{
				Msg: "标签不存在",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, model.MsgResponse{
				Msg: "服务器错误",
			})
			return
		}
	}

	c.JSON(http.StatusOK, model.CreateTagResponseSuccessBody{Resource: model.Tag{
		ID:        r.ID,
		UserID:    r.UserID,
		Name:      r.Name,
		Sign:      r.Sign,
		Kind:      r.Kind,
		DeletedAt: r.DeletedAt,
	}})
}

func (ctl *TagController) Destroy(c *gin.Context) {
	// TODO implement me
	panic("implement me")
}
