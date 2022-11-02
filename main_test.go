package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := SetUpRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/healthchecker", nil)

	if err != nil {
		t.FailNow()
	}
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"message\":\"server connected\"}", w.Body.String())
}
