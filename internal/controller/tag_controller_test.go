package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
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
		Name: null.String{NullString: sql.NullString{
			String: "test 333",
			Valid:  true,
		}},
		Sign: null.String{NullString: sql.NullString{
			String: "🚗",
			Valid:  true,
		}},
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
