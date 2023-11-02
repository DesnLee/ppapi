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

// 验证账单金额
func validateAmount(amount int64) error {
    if amount <= 0 {
        return fmt.Errorf("金额必须大于 0")
    }
    return nil
}

// 验证账单类型
func validateKind(kind string) error {
    if kind != "income" && kind != "expenses" {
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
    if parsed, err := time.Parse("2006-01-02", happenedAt); err != nil {
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
    tags, err := database.Q.FindTagsByIDs(database.DBCtx, ids)
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

// 验证页码
func validatePage(page int64) error {
    if page <= 0 {
        return fmt.Errorf("page 必须大于 0")
    }
    return nil
}

// 验证每页条目数
func validateSize(size int64) error {
    if size <= 0 {
        return fmt.Errorf("size 必须大于 0")
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
        parsedAfter, _ := time.Parse("2006-01-02", after)
        parsedBefore, _ := time.Parse("2006-01-02", before)

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
    if err := validateKind(string(b.Kind)); err != nil {
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
        HappenedAt: r.HappenedAt.Time.Format("2006-01-02"),
        TagIDs:     b.TagIDs,
    }, nil
}

// 验证获取账单列表请求体
func validateGetItemsRequestBody(b *model.GetItemsRequestBody) error {
    if err := validatePage(b.Page); err != nil {
        return err
    }
    if err := validateSize(b.Size); err != nil {
        return err
    }
    if err := validateDateRange(b.HappenedAfter, b.HappenedBefore); err != nil {
        return err
    }

    return nil
}

// GetAndCountItemsByUserID 获取账单列表
func GetAndCountItemsByUserID(uid pgtype.UUID, b model.GetItemsRequestBody) (model.GetItemsResponseSuccessBody, error) {
    var res model.GetItemsResponseSuccessBody

    tx, err := database.DB.Begin(database.DBCtx)
    if err != nil {
        log.Println("[Create Database Transaction Failed]: ", err)
        return res, fmt.Errorf("服务器错误")
    }
    defer func(tx *pgx.Tx) {
        _ = (*tx).Rollback(database.DBCtx)
    }(&tx)

    qtx := database.Q.WithTx(tx)

    // 初始化未传字段
    if b.Page == 0 {
        b.Page = 1
    }
    if b.Size == 0 {
        b.Size = 10
    }
    if err := validateGetItemsRequestBody(&b); err != nil {
        return res, err
    }

    if b.HappenedBefore == "" {
        b.HappenedBefore = time.Now().Format("2006-01-02")
    }
    before, _ := pkg.CreatePgTimeTZ(b.HappenedBefore)

    if b.HappenedAfter == "" {
        b.HappenedAfter = "1970-01-01"
    }
    after, _ := pkg.CreatePgTimeTZ(b.HappenedAfter)

    items, err := qtx.ListItemsByUserID(database.DBCtx, sqlcExec.ListItemsByUserIDParams{
        UserID:         uid,
        Offset:         int32((b.Page - 1) * b.Size),
        Limit:          int32(b.Size),
        HappenedAfter:  after,
        HappenedBefore: before,
    })
    if err != nil {
        return res, fmt.Errorf("服务器错误")
    }

    if len(items) == 0 {
        return res, nil
    }
    for _, item := range items {
        res.Resources = append(res.Resources, model.GetItemsResponseData{
            CreateItemResponseData: model.CreateItemResponseData{
                ID:         item.ID,
                Amount:     item.Amount,
                Kind:       item.Kind,
                HappenedAt: item.HappenedAt.Time.Format("2006-01-02"),
                TagIDs:     item.TagIds,
            },
            // Tags: []model.TagResponse{},
        })
        // for _, tag := range item.Tags {
        //     tagRow := tag.(sqlcExec.Tag)
        //     res.Resources[i].TagIDs = append(res.Resources[i].TagIDs, tagRow.ID)
        //     res.Resources[i].Tags = append(res.Resources[i].Tags, model.TagResponse{
        //         ID:   tagRow.ID,
        //         Name: tagRow.Name,
        //         Sign: tagRow.Sign,
        //         Kind: tagRow.Kind,
        //     })
        // }
    }

    count, err := qtx.CountItemsByUserID(database.DBCtx, uid)
    if err != nil {
        return res, fmt.Errorf("服务器错误")
    }
    res.Pager = model.Pager{
        Page:    b.Page,
        PerPage: b.Size,
        Count:   count,
    }
    return res, nil
}
