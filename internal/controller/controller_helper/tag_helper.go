package controller_helper

import (
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/constants"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/model"
)

// ValidateCreateTagRequestBody 验证创建标签请求体
func ValidateCreateTagRequestBody(b *model.CreateTagRequestBody) error {
	if err := validateKind(b.Kind); err != nil {
		return err
	}
	if b.Name == "" {
		return fmt.Errorf("标签名称不能为空")
	}
	if b.Sign == "" {
		return fmt.Errorf("标签图标不能为空")
	}
	return nil
}

// ValidateUpdateTagRequestBody 验证更新标签请求体
func ValidateUpdateTagRequestBody(b *model.UpdateTagRequestBody) error {
	if b.Kind.Valid && b.Kind.String != "" {
		if err := validateKind(b.Kind.String); err != nil {
			return err
		}
	}
	return nil
}

// 生成查询条件
func generateQueryTagsCondition(uid pgtype.UUID, b model.GetTagsRequestBody) (sqlcExec.ListTagsByUserIDWithConditionParams, error) {
	condition := sqlcExec.ListTagsByUserIDWithConditionParams{
		UserID: uid,
		Kind:   b.Kind,
	}
	// 初始化未传字段
	if b.Size == 0 {
		condition.Limit = 10
		b.Size = 10
	} else {
		condition.Limit = int32(b.Size)
	}
	if b.Page == 0 {
		condition.Offset = 0
		b.Page = 1
	} else {
		condition.Offset = int32((b.Page - 1) * b.Size)
	}

	return condition, nil
}

// GetAndCountTagsByUserID 获取标签列表
func GetAndCountTagsByUserID(uid pgtype.UUID, b model.GetTagsRequestBody) (model.GetTagsResponseSuccessBody, error) {
	res := model.GetTagsResponseSuccessBody{}

	tx, err := database.DB.Begin(database.DBCtx)
	if err != nil {
		log.Println("[Create Database Transaction Failed]: ", err)
		return res, fmt.Errorf("服务器错误")
	}
	defer func(tx *pgx.Tx) {
		_ = (*tx).Rollback(database.DBCtx)
	}(&tx)

	qtx := database.Q.WithTx(tx)

	if b.Kind == "" {
		b.Kind = constants.KindExpenses
	}
	if err := validateKind(b.Kind); err != nil {
		return res, err
	}

	condition, err := generateQueryTagsCondition(uid, b)
	if err != nil {
		return res, err
	}
	rows, err := qtx.ListTagsByUserIDWithCondition(database.DBCtx, condition)
	if err != nil {
		return res, fmt.Errorf("服务器错误")
	}

	count, err := qtx.CountTagsByUserIDWithCondition(database.DBCtx, sqlcExec.CountTagsByUserIDWithConditionParams{
		UserID: uid,
		Kind:   b.Kind,
	})
	if err != nil {
		return res, fmt.Errorf("服务器错误")
	}

	for _, row := range rows {
		res.Resources = append(res.Resources, model.Tag{
			ID:        row.ID,
			UserID:    row.UserID,
			Name:      row.Name,
			Sign:      row.Sign,
			Kind:      row.Kind,
			DeletedAt: row.DeletedAt,
		})
	}

	res.Pager = model.Pager{
		Page:    int64(condition.Offset/condition.Limit + 1),
		PerPage: int64(condition.Limit),
		Count:   count,
	}
	return res, nil
}
