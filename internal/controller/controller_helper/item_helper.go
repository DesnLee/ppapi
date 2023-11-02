package controller_helper

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

func validateAmount(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("金额必须大于 0")
	}
	return nil
}

func validateKind(kind string) error {
	if kind != "income" && kind != "expenses" {
		return fmt.Errorf("kind 必须为 income 或 expenses")
	}
	return nil
}

func validateHappenedAt(happenedAt string) error {
	if happenedAt == "" {
		return fmt.Errorf("账单时间不能为空")
	}

	now := time.Now()
	if parsed, err := time.Parse("2006-01-02", happenedAt); err != nil {
		return fmt.Errorf("账单时间格式错误")
	} else if parsed.After(now) {
		return fmt.Errorf("账单时间不能晚于当前时间")
	}
	return nil
}

func validateTagIDs(ids []int64) error {
	if len(ids) == 0 {
		return fmt.Errorf("至少选择一个标签")
	}

	// 创建一个map，用于存储切片中的元素和它们的出现次数
	seen := make(map[int64]bool)

	// 遍历切片
	for _, id := range ids {
		// 如果元素已经在map中存在，表示有重复项
		if seen[id] {
			return fmt.Errorf("标签不能重复")
		}
		// 否则，将元素添加到map中
		seen[id] = true
	}
	// 如果没有找到重复项，则去数据库中查询是否存在这些标签
	tags, err := database.Q.FindTagsByIDs(database.DBCtx, ids)
	if err != nil {
		return fmt.Errorf("服务器错误")
	}
	if len(tags) != len(ids) {
		return fmt.Errorf("存在无效标签")
	}

	return nil
}

func ValidateCreateItemRequestBody(b *model.CreateItemRequestBody) error {
	if err := validateAmount(b.Amount); err != nil {
		return err
	}
	if err := validateKind(string(b.Kind)); err != nil {
		return err
	}
	if err := validateHappenedAt(b.HappenedAt); err != nil {
		return err
	}
	if err := validateTagIDs(b.TagIDs); err != nil {
		return err
	}
	return nil
}

func CreateItem(uid pgtype.UUID, b *model.CreateItemRequestBody) (model.CreateItemResponseData, error) {
	var item model.CreateItemResponseData

	tx, err := database.DB.Begin(database.DBCtx)
	if err != nil {
		log.Println("[Create Database Transaction Failed]: ", err)
		return item, fmt.Errorf("服务器错误")
	}
	defer func(tx *pgx.Tx) {
		_ = (*tx).Rollback(database.DBCtx)
	}(&tx)

	qtx := database.Q.WithTx(tx)

	tm, err := pkg.CreatePgTimeTZ(b.HappenedAt)
	if err != nil {
		return item, fmt.Errorf("happenedAt 格式错误")
	}

	r, err := qtx.CreateItem(database.DBCtx, sqlcExec.CreateItemParams{
		UserID:     uid,
		Amount:     b.Amount,
		Kind:       b.Kind,
		HappenedAt: tm,
	})

	if err != nil {
		log.Println("[Create Item Failed]: ", err)
		return item, fmt.Errorf("服务器错误")
	}

	// 创建账单标签关联
	var relations []sqlcExec.CreateItemTagRelationsParams
	for _, id := range b.TagIDs {
		relations = append(relations, sqlcExec.CreateItemTagRelationsParams{
			ItemID: r.ID,
			TagID:  id,
		})
	}
	if _, err = qtx.CreateItemTagRelations(database.DBCtx, relations); err != nil {
		log.Println("[Create ItemTagRelations Failed]: ", err)
		return item, fmt.Errorf("服务器错误")
	}

	err = tx.Commit(database.DBCtx)
	if err != nil {
		return item, fmt.Errorf("服务器错误")
	}
	return model.CreateItemResponseData{
		ID:         r.ID,
		Amount:     r.Amount,
		Kind:       r.Kind,
		HappenedAt: r.HappenedAt.Time.Format("2006-01-02"),
		TagIDs:     b.TagIDs,
	}, nil
}
