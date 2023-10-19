package controller_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/db/sqlcExec"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/router"
	"ppapi.desnlee.com/pkg"
)

func TestSession(t *testing.T) {
	r := router.New()
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

	body := gin.H{
		"email": email,
		"code":  code,
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/session", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	// 校验返回的消息
	schema := struct {
		JWT string `json:"jwt"`
	}{}
	if err = json.Unmarshal(w.Body.Bytes(), &schema); err != nil {
		t.Error("Unmarshal Response Body Error: ", err)
	}

	fmt.Println("JWT: ", schema.JWT)

	assert.Equal(t, http.StatusOK, w.Code)
}
