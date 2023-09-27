package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/config"
	"ppapi.desnlee.com/internal/database"
	"ppapi.desnlee.com/internal/router"
)

func init() {
	database.Connect()

	config.LoadConfig()
	viper.Set("EMAIL.SMTP.HOST", "localhost")
	viper.Set("EMAIL.SMTP.PORT", "1025")
}

func TestValidationCode(t *testing.T) {
	r := router.New()
	w := httptest.NewRecorder()

	oldCount, _ := database.Q.CountValidationCodes(database.DBCtx)

	req, _ := http.NewRequest("POST", "/api/v1/validation_code", strings.NewReader(`{"email":"test@qq.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	newCount, _ := database.Q.CountValidationCodes(database.DBCtx)

	assert.Equal(t, oldCount+1, newCount)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
