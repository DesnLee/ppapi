package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

func TestCreateTag(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&TagController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 发送请求
	body := model.CreateTagRequestBody{
		Name: "test",
		Sign: "😄",
		Kind: sqlcExec.KindExpenses,
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/tags", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.CreateTagResponseSuccessBody{}
	if err := json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test", schema.Resource.Name)
	assert.Equal(t, "😄", schema.Resource.Sign)
	assert.Equal(t, sqlcExec.KindExpenses, schema.Resource.Kind)
	assert.Equal(t, u.ID, schema.Resource.UserID)
	assert.Equal(t, pgtype.Timestamptz{}, schema.Resource.DeletedAt)
}

func TestGetTag(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&TagController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 创建标签
	name := "test"
	sign := "😄"
	kind := sqlcExec.KindExpenses
	tag, _ := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   &name,
		Sign:   &sign,
		Kind:   &kind,
	})

	// 发送请求
	url := fmt.Sprintf("/api/v1/tags/%d", tag.ID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.GetTagResponseSuccessBody{}
	if err := json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, *tag.Kind, schema.Resource.Kind)
	assert.Equal(t, *tag.Name, schema.Resource.Name)
	assert.Equal(t, *tag.Sign, schema.Resource.Sign)
	assert.Equal(t, tag.UserID, schema.Resource.UserID)
	assert.Equal(t, tag.DeletedAt, schema.Resource.DeletedAt)
}

func TestUpdateTag(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&TagController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 创建标签
	name := "test"
	sign := "😄"
	kind := sqlcExec.KindExpenses
	tag, _ := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   &name,
		Sign:   &sign,
		Kind:   &kind,
	})

	// 发送请求
	body := model.UpdateTagRequestBody{
		Name: "test 2",
		Sign: "🚗",
	}
	bodyStr, _ := json.Marshal(body)
	url := fmt.Sprintf("/api/v1/tags/%d", tag.ID)
	req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.UpdateTagResponseSuccessBody{}
	if err := json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, *tag.Kind, schema.Resource.Kind)
	assert.Equal(t, body.Name, schema.Resource.Name)
	assert.Equal(t, body.Sign, schema.Resource.Sign)
	assert.Equal(t, tag.UserID, schema.Resource.UserID)
	assert.Equal(t, tag.DeletedAt, schema.Resource.DeletedAt)
}
