package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"ppapi.desnlee.com/internal/router"
)

func TestPing(t *testing.T) {
	r := router.New()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"msg\":\"pong\"}", w.Body.String())
}
