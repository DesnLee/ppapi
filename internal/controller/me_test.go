package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/jwt_helper"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

func TestMe(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&MeController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/me", strings.NewReader(""))
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.MsgResponse{}
	if err := json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMeWithJWT(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	(&MeController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	// 创建用户
	u, err := database.Q.CreateUser(database.DBCtx, "test@qq.com")
	if err != nil {
		log.Fatal("Create User Error: ", err)
	}

	// 生成 JWT
	jwtStr, err := jwt_helper.GenerateJWT(u.ID)
	if err != nil {
		log.Fatal("Generate JWT Error: ", err)
	}

	// 发送请求
	req, _ := http.NewRequest("GET", "/api/v1/me", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtStr)
	r.ServeHTTP(w, req)

	// 反序列化返回的消息
	schema := model.ResourceResponse[model.MeResponseBody]{
		Resource: model.MeResponseBody{},
	}
	if err = json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, u.ID, schema.Resource.ID)
}
