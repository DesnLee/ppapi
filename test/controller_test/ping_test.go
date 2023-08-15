package controller_test

import (
	"net/http"
	"net/http/httptest"
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

func TestPing(t *testing.T) {
	r := router.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"msg\":\"pong\"}", w.Body.String())
}
