package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/pkg"
)

func TestPing(t *testing.T) {
	r, cleaner := pkg.InitTestEnv()
	defer cleaner()
	r.GET("/ping", PingHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"msg\":\"pong\"}", w.Body.String())
}
