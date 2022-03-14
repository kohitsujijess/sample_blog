package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kohitsujijess/sample_blog/router"

	"github.com/stretchr/testify/assert"
)

func TestController(t *testing.T) {
	router := router.Init()

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Entries of sample blog", rec.Body.String())
}
