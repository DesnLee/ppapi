package controller_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/internal/router"
)

func TestValidationCode(t *testing.T) {
	r := router.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/validation_code", strings.NewReader(`{"email":"test@qq.com"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
