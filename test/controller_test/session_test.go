package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/internal/router"
)

func TestSession(t *testing.T) {
	r := router.New()
	w := httptest.NewRecorder()

	body := gin.H{
		"email": "",
		"code":  "",
	}
	bodyStr, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/session", strings.NewReader(string(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
