package controller_helper

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/constants"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

// 验证账单金额
func validateAmount(amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("金额必须大于 0")
	}
	return nil
}

// 验证账单类型
func validateKind(kind string) error {
	if kind != constants.KindIncome && kind != constants.KindExpenses {
		return fmt.Errorf("kind 必须为 income 或 expenses")
	}
	return nil
}

// 验证账单发生时间
func validateHappenedAt(happenedAt string) error {
	if happenedAt == "" {
		return fmt.Errorf("账单时间不能为空")
	}

	now := time.Now()
	if parsed, err := time.Parse(time.RFC3339, happenedAt); err != nil {
		return fmt.Errorf("账单时间格式错误")
	} else if parsed.After(now) {
		return fmt.Errorf("账单时间不能晚于当前时间")
	}
	return nil
}

// 验证账单标签
func validateTagIDs(uid pgtype.UUID, ids []int64) error {
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
	tags, err := database.Q.FindTagsByIDs(database.DBCtx, sqlcExec.FindTagsByIDsParams{
		TagIds: ids,
		UserID: uid,
	})
	if err != nil {
		return fmt.Errorf("服务器错误")
	}
	if len(tags) != len(ids) {
		return fmt.Errorf("存在无效标签")
	}
	for _, tag := range tags {
		if tag.UserID != uid {
			return fmt.Errorf("存在无效标签")
		}
	}

	return nil
}

// 验证日期范围
func validateDateRange(after, before string) error {
	if before != "" {
		err := validateHappenedAt(before)
		if err != nil {
			return err
		}
	}
	if after != "" {
		err := validateHappenedAt(after)
		if err != nil {
			return err
		}
	}
	if after != "" && before != "" {
		parsedAfter, _ := time.Parse(time.RFC3339, after)
		parsedBefore, _ := time.Parse(time.RFC3339, before)

		if parsedAfter.After(parsedBefore) {
			return fmt.Errorf("after 不能晚于 before")
		}
	}

	return nil
}

// ValidateCreateItemRequestBody 验证创建账单请求体
func ValidateCreateItemRequestBody(uid pgtype.UUID, b *model.CreateItemRequestBody) error {
	if err := validateAmount(b.Amount); err != nil {
		return err
	}
	if err := validateKind(b.Kind); err != nil {
		return err
	}
	if err := validateHappenedAt(b.HappenedAt); err != nil {
		return err
	}
	if err := validateTagIDs(uid, b.TagIDs); err != nil {
		return err
	}
	return nil
}

// CreateItem 创建账单
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
		HappenedAt: r.HappenedAt,
		TagIDs:     b.TagIDs,
	}, nil
}

// 生成时间区间
func generateDateRange(a string, b string) (pgtype.Timestamptz, pgtype.Timestamptz) {
	if a == "" {
		// RFC3339 的最小值
		a = time.Now().AddDate(-50, 0, 0).Local().Format(time.RFC3339)
	}
	after, _ := pkg.CreatePgTimeTZ(a)

	if b == "" {
		b = time.Now().AddDate(0, 0, 1).Local().Format(time.RFC3339)
	}
	before, _ := pkg.CreatePgTimeTZ(b)

	return after, before
}

// 生成查询条件
func generateQueryItemsCondition(uid pgtype.UUID, b model.GetItemsRequestBody) (sqlcExec.ListItemsByUserIDWithConditionParams, error) {
	condition := sqlcExec.ListItemsByUserIDWithConditionParams{
		UserID: uid,
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

	after, before := generateDateRange(b.HappenedAfter, b.HappenedBefore)
	condition.HappenedBefore = before
	condition.HappenedAfter = after

	return condition, nil
}

// GetAndCountItemsByUserID 获取账单列表
func GetAndCountItemsByUserID(uid pgtype.UUID, b model.GetItemsRequestBody) (model.GetItemsResponseSuccessBody, error) {
	res := model.GetItemsResponseSuccessBody{
		Resources: []model.GetItemsResponseData{},
	}

	tx, err := database.DB.Begin(database.DBCtx)
	if err != nil {
		log.Println("[Create Database Transaction Failed]: ", err)
		return res, fmt.Errorf("服务器错误")
	}
	defer func(tx *pgx.Tx) {
		_ = (*tx).Rollback(database.DBCtx)
	}(&tx)

	qtx := database.Q.WithTx(tx)

	if err := validateDateRange(b.HappenedAfter, b.HappenedBefore); err != nil {
		return res, err
	}

	condition, err := generateQueryItemsCondition(uid, b)
	if err != nil {
		return res, err
	}
	rows, err := qtx.ListItemsByUserIDWithCondition(database.DBCtx, condition)
	if err != nil {
		return res, fmt.Errorf("服务器错误")
	}

	count, err := qtx.CountItemsByUserIDWithCondition(database.DBCtx, sqlcExec.CountItemsByUserIDWithConditionParams{
		UserID:         uid,
		HappenedAfter:  condition.HappenedAfter,
		HappenedBefore: condition.HappenedBefore,
	})
	if err != nil {
		return res, fmt.Errorf("服务器错误")
	}

	for _, item := range rows {
		res.Resources = append(res.Resources, model.GetItemsResponseData{
			CreateItemResponseData: model.CreateItemResponseData{
				ID:         item.ID,
				Amount:     item.Amount,
				Kind:       item.Kind,
				HappenedAt: item.HappenedAt,
				TagIDs:     item.TagIds,
			},
			Tags: []model.Tag{},
		})
	}

	// 创建 map
	tagIDsMap := make(map[int64]bool)
	for _, item := range rows {
		for _, tagID := range item.TagIds {
			tagIDsMap[tagID] = true
		}
	}
	// 提取 map 中的 key
	tagIDs := make([]int64, 0, len(tagIDsMap))
	for tagID := range tagIDsMap {
		tagIDs = append(tagIDs, tagID)
	}
	// 查询标签
	tags, err := qtx.FindTagsByIDs(database.DBCtx, sqlcExec.FindTagsByIDsParams{
		TagIds: tagIDs,
		UserID: uid,
	})
	if err != nil {
		return res, fmt.Errorf("服务器错误")
	}

	// 将 tags 详情 插入到 res 的每个 Resoure 的 tags 字段中
	for i, item := range res.Resources {
		for _, id := range item.TagIDs {
			for _, tag := range tags {
				if id == tag.ID {
					res.Resources[i].Tags = append(res.Resources[i].Tags, model.Tag{
						ID:   tag.ID,
						Name: tag.Name,
						Sign: tag.Sign,
					})
				}
			}
		}
	}

	res.Pager = model.Pager{
		Page:    int64(condition.Offset/condition.Limit + 1),
		PerPage: int64(condition.Limit),
		Count:   count,
	}
	return res, nil
}

// ValidateGetBalanceRequestBody 验证收支信息请求体
func ValidateGetBalanceRequestBody(b *model.GetBalanceRequestBody) error {
	return validateDateRange(b.HappenedAfter, b.HappenedBefore)
}

// GetBalance 获取收支信息
func GetBalance(uid pgtype.UUID, b model.GetBalanceRequestBody) (model.GetBalanceResponseData, error) {
	var res model.GetBalanceResponseData
	after, before := generateDateRange(b.HappenedAfter, b.HappenedBefore)
	items, err := database.Q.ListItemsByUserIDWithDateRange(database.DBCtx, sqlcExec.ListItemsByUserIDWithDateRangeParams{
		UserID:         uid,
		HappenedAfter:  after,
		HappenedBefore: before,
	})
	if err != nil {
		fmt.Println("[List Items By UserID With Date Range Failed]: ", err)
		return res, fmt.Errorf("服务器错误")
	}

	for _, item := range items {
		if item.Kind == constants.KindIncome {
			res.Income += item.Amount
		} else if item.Kind == constants.KindExpenses {
			res.Expenses += item.Amount
		}
	}
	res.Balance = res.Income - res.Expenses
	return res, nil
}
