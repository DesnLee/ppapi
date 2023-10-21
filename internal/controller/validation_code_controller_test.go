package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/model"
	"ppapi.desnlee.com/pkg"
)

func TestValidationCode(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	viper.Set("EMAIL.SMTP.HOST", "localhost")
	viper.Set("EMAIL.SMTP.PORT", "1025")
	(&ValidationCodeController{}).Register(r.Group("/api"))

	w := httptest.NewRecorder()

	oldCount, _ := database.Q.CountValidationCodes(database.DBCtx)

	body := model.ValidationCodeRequestBody{
		Email: "test@qq.com",
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/validation_code", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	newCount, _ := database.Q.CountValidationCodes(database.DBCtx)

	assert.Equal(t, oldCount+1, newCount)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
