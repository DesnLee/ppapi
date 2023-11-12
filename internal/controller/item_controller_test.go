package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

func TestCreateItem(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&ItemController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 创建标签
	name := "test"
	sign := "😄"
	kind := sqlcExec.KindExpenses
	tag, err := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   &name,
		Sign:   &sign,
		Kind:   &kind,
	})
	if err != nil {
		log.Fatal("Create Tag Error: ", err)
	}

	// 发送请求
	body := model.CreateItemRequestBody{
		Amount:     1999,
		Kind:       sqlcExec.KindExpenses,
		HappenedAt: time.Now().Local().Format(time.RFC3339),
		TagIDs:     []int64{tag.ID},
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/items", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.ResourceResponse[model.CreateItemResponseData]{
		Resource: model.CreateItemResponseData{},
	}
	if err = json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	item, err := database.Q.FindItemByID(database.DBCtx, schema.Resource.ID)
	if err != nil {
		t.Error("Find Item By ID Error: ", err)
	}
	tagIDs, err := database.Q.FindTagIDsByItemID(database.DBCtx, schema.Resource.ID)
	if err != nil {
		t.Error("Find Tag IDs By Item ID Error: ", err)
	}
	if len(tagIDs) != 1 {
		t.Error("Find Tag IDs By Item ID Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(1999), item.Amount)
	assert.Equal(t, sqlcExec.KindExpenses, string(item.Kind))
	assert.Equal(t, body.HappenedAt, item.HappenedAt.Time.Format(time.RFC3339))
	assert.Equal(t, []int64{tag.ID}, tagIDs)
}

func TestGetMultiItems(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&ItemController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 创建标签
	name := "test"
	sign := "😄"
	kind := sqlcExec.KindExpenses
	tag, err := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   &name,
		Sign:   &sign,
		Kind:   &kind,
	})
	if err != nil {
		log.Fatal("Create Tag Error: ", err)
	}

	// 创建账目
	for i := 0; i < 20; i++ {
		tm, err := pkg.CreatePgTimeTZ(time.Now().Local().Format(time.RFC3339))
		if err != nil {
			log.Fatal("Create Time Error: ", err)
		}
		item, err := database.Q.CreateItem(database.DBCtx, sqlcExec.CreateItemParams{
			UserID:     u.ID,
			Amount:     1999,
			Kind:       sqlcExec.KindExpenses,
			HappenedAt: tm,
		})
		if err != nil {
			log.Fatal("Create Item Error: ", err)
		}
		_, err = database.Q.CreateItemTagRelations(database.DBCtx, []sqlcExec.CreateItemTagRelationsParams{
			{
				ItemID: item.ID,
				TagID:  tag.ID,
			},
		})
	}

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/items", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.ResourcesResponse[model.GetItemsResponseData]{}
	if err = json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 10, len(schema.Resources))
	assert.Equal(t, int64(1999), schema.Resources[0].Amount)
	assert.Equal(t, int64(1), schema.Pager.Page)
	assert.Equal(t, int64(20), schema.Pager.Count)
}

func TestGetBalance(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&ItemController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 创建标签
	name := "test"
	sign := "😄"
	kind := sqlcExec.KindExpenses
	tag, err := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   &name,
		Sign:   &sign,
		Kind:   &kind,
	})
	if err != nil {
		log.Fatal("Create Tag Error: ", err)
	}

	// 创建第一批
	for i := 0; i < 10; i++ {
		tm, _ := pkg.CreatePgTimeTZ(time.Date(2022, 12, 31, 23, 59, 59, 0, time.Local).Format(time.RFC3339))
		item, err := database.Q.CreateItem(database.DBCtx, sqlcExec.CreateItemParams{
			UserID:     u.ID,
			Amount:     2000,
			Kind:       sqlcExec.KindExpenses,
			HappenedAt: tm,
		})
		if err != nil {
			log.Fatal("Create Item Error: ", err)
		}
		_, err = database.Q.CreateItemTagRelations(database.DBCtx, []sqlcExec.CreateItemTagRelationsParams{
			{
				ItemID: item.ID,
				TagID:  tag.ID,
			},
		})
	}
	// 创建第二批
	for i := 0; i < 10; i++ {
		tm, _ := pkg.CreatePgTimeTZ(time.Date(2023, 01, 01, 12, 00, 00, 0, time.Local).Format(time.RFC3339))
		item, err := database.Q.CreateItem(database.DBCtx, sqlcExec.CreateItemParams{
			UserID:     u.ID,
			Amount:     2000,
			Kind:       sqlcExec.KindExpenses,
			HappenedAt: tm,
		})
		if err != nil {
			log.Fatal("Create Item Error: ", err)
		}
		_, err = database.Q.CreateItemTagRelations(database.DBCtx, []sqlcExec.CreateItemTagRelationsParams{
			{
				ItemID: item.ID,
				TagID:  tag.ID,
			},
		})
	}
	// 创建第三批
	for i := 0; i < 10; i++ {
		tm, _ := pkg.CreatePgTimeTZ(time.Date(2023, 01, 10, 12, 00, 00, 0, time.Local).Format(time.RFC3339))
		item, err := database.Q.CreateItem(database.DBCtx, sqlcExec.CreateItemParams{
			UserID:     u.ID,
			Amount:     2000,
			Kind:       sqlcExec.KindExpenses,
			HappenedAt: tm,
		})
		if err != nil {
			log.Fatal("Create Item Error: ", err)
		}
		_, err = database.Q.CreateItemTagRelations(database.DBCtx, []sqlcExec.CreateItemTagRelationsParams{
			{
				ItemID: item.ID,
				TagID:  tag.ID,
			},
		})
	}

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/items/balance?happened_after=2023-01-01T00:00:00%2B08:00&happened_before=2023-01-02T00:00:00%2B08:00",
		nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.GetBalanceResponseData{}
	if err = json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(20000), schema.Expenses)
	assert.Equal(t, int64(0), schema.Income)
	assert.Equal(t, int64(-20000), schema.Balance)
}
