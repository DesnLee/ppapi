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

func TestSession(t *testing.T) {
	r := pkg.SetupTest()
	(&SessionController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	email := "test@qq.com"
	code, err := pkg.GenerateRandomCode(6)

	if err != nil {
		log.Fatal("Create Random Code Error: ", err)
	}

	_, err = database.Q.CreateValidationCode(database.DBCtx, sqlcExec.CreateValidationCodeParams{
		Email: email,
		Code:  code,
	})

	if err != nil {
		log.Fatal("Insert Code to DB Error: ", err)
	}

	body := model.SessionRequestBody{
		Email: email,
		Code:  code,
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/session", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// 校验返回的消息
	schema := model.SessionResponseBody{}
	if err = json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	u, err := database.Q.FindUserByEmail(database.DBCtx, email)
	if err != nil {
		t.Error("Find User By Email Error: ", err)
	}

	uid, err := jwt_helper.GetJWTUserID(schema.JWT)
	if err != nil {
		t.Error("Get JWT User ID Error: ", err)
	}

	assert.Equal(t, u.ID, uid)
	assert.Equal(t, http.StatusOK, w.Code)
}
