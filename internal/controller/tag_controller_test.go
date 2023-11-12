package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		Kind: "expenses",
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/tags", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.ResourceResponse[model.Tag]{
		Resource: model.Tag{},
	}
	if err := json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	// 校验返回的数据
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test", schema.Resource.Name)
	assert.Equal(t, "😄", schema.Resource.Sign)
	assert.Equal(t, "expenses", string(schema.Resource.Kind))
	assert.Equal(t, u.ID, schema.Resource.UserID)
}
