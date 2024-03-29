package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/constants"
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
		Kind: constants.KindExpenses,
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
	assert.Equal(t, constants.KindExpenses, schema.Resource.Kind)
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
	tag, _ := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   "test",
		Sign:   "😄",
		Kind:   constants.KindExpenses,
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
	assert.Equal(t, tag.Kind, schema.Resource.Kind)
	assert.Equal(t, tag.Name, schema.Resource.Name)
	assert.Equal(t, tag.Sign, schema.Resource.Sign)
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
	tag, _ := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   "test",
		Sign:   "😄",
		Kind:   constants.KindExpenses,
	})

	// 发送请求
	body := model.UpdateTagRequestBody{
		Name: pgtype.Text{String: "test 333", Valid: true},
		Sign: pgtype.Text{String: "🚗", Valid: true},
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
	assert.Equal(t, tag.Kind, schema.Resource.Kind)
	assert.Equal(t, body.Name.String, schema.Resource.Name)
	assert.Equal(t, body.Sign.String, schema.Resource.Sign)
	assert.Equal(t, tag.UserID, schema.Resource.UserID)
	assert.Equal(t, tag.DeletedAt, schema.Resource.DeletedAt)
}

func TestDeleteTag(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&TagController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 创建标签
	tag, _ := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
		UserID: u.ID,
		Name:   "test",
		Sign:   "😄",
		Kind:   constants.KindExpenses,
	})

	// 发送请求
	url := fmt.Sprintf("/api/v1/tags/%d", tag.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	_, err := database.Q.FindTagByID(database.DBCtx, sqlcExec.FindTagByIDParams{
		ID:     tag.ID,
		UserID: u.ID,
	})
	assert.Equal(t, pgx.ErrNoRows, err)
}

func TestGetMultiTags(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&TagController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	u, jwt := pkg.TestCreateUserAndJWT()

	// 创建标签
	for i := 0; i < 20; i++ {
		_, err := database.Q.CreateTag(database.DBCtx, sqlcExec.CreateTagParams{
			UserID: u.ID,
			Name:   fmt.Sprintf("test %d", i),
			Sign:   "😄",
			Kind:   constants.KindExpenses,
		})
		if err != nil {
			t.Error("Create Tag Error: ", err)
		}
	}

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/tags", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.GetTagsResponseSuccessBody{}
	if err := json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 10, len(schema.Resources))
	assert.Equal(t, "test 19", schema.Resources[0].Name)
	assert.Equal(t, int64(1), schema.Pager.Page)
	assert.Equal(t, int64(20), schema.Pager.Count)
}
