package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestControllerHealth(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sonarqube-to-gitlab-webhook/v1/health")

	// Assertions
	if assert.NoError(t, CheckHealth(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}