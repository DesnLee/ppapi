package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/jwt_helper"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

func TestCreateItem(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&ItemController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	// 创建用户
	u, err := database.Q.CreateUser(database.DBCtx, "test@qq.com")
	if err != nil {
		log.Fatal("Create User Error: ", err)
	}
	// 创建标签
	tag, err := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   "test",
		Sign:   "😄",
		Kind:   "expenses",
	})
	if err != nil {
		log.Fatal("Create Tag Error: ", err)
	}

	// 生成 JWT
	jwtStr, err := jwt_helper.GenerateJWT(u.ID)
	if err != nil {
		log.Fatal("Generate JWT Error: ", err)
	}

	// 发送请求
	body := model.CreateItemRequestBody{
		Amount:     1999,
		Kind:       "expenses",
		HappenedAt: "2023-01-01",
		TagIDs:     []int64{tag.ID},
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/item", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtStr)
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
	assert.Equal(t, "expenses", string(item.Kind))
	assert.Equal(t, "2023-01-01", item.HappenedAt.Time.Format("2006-01-02"))
	assert.Equal(t, []int64{tag.ID}, tagIDs)
}
