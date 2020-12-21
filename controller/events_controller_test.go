package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/appcontext"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	localtest "github.com/betorvs/sonarqube-to-gitlab-webhook/test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPostEvents(t *testing.T) {
	// Setup
	e := echo.New()
	config.Values.SonarqubeSecret = "Absent"
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/sonarqube-to-gitlab-webhook/v1/events")

	// Assertions
	if assert.NoError(t, ReceiveEvents(c)) {
		assert.Equal(t, http.StatusNotImplemented, rec.Code)
	}

	e1 := echo.New()
	config.Values.SonarqubeSecret = "6885df05d1989243b9978b3a32b626b7a2e3091554b040f4"
	req1 := httptest.NewRequest(http.MethodPost, "/", nil)
	rec1 := httptest.NewRecorder()
	c1 := e1.NewContext(req1, rec1)
	c1.SetPath("/sonarqube-to-gitlab-webhook/v1/events")

	// Assertions
	if assert.NoError(t, ReceiveEvents(c1)) {
		assert.Equal(t, http.StatusForbidden, rec1.Code)
	}

	e2 := echo.New()
	header2 := "0488c4999405ba55e256a2b412695526f98b88435c5a7fe2484d0e19421aa5b8"
	req2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(localtest.TestBadJSON))
	req2.Header.Add("X-Sonar-Webhook-Hmac-Sha256", header2)
	rec2 := httptest.NewRecorder()
	c2 := e2.NewContext(req2, rec2)
	c2.SetPath("/sonarqube-to-gitlab-webhook/v1/events")

	// Assertions
	if assert.NoError(t, ReceiveEvents(c2)) {
		assert.Equal(t, http.StatusBadRequest, rec2.Code)
	}

	e3 := echo.New()
	bodyBytes := []byte(localtest.TestValidJSON)
	bodyString := string(bodyBytes)
	header3 := "74d2b38fef2c09e9b02156d27bb23f76d22ab505dc1ad8d8881e6be01fe0b41b"
	req3 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyString))
	req3.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req3.Header.Add("X-Sonar-Webhook-Hmac-Sha256", header3)
	rec3 := httptest.NewRecorder()
	c3 := e3.NewContext(req3, rec3)
	c3.SetPath("/sonarqube-to-gitlab-webhook/v1/events")
	appcontext.Current.Add(appcontext.Logger, localtest.InitMockLogger)
	// Assertions
	if assert.NoError(t, ReceiveEvents(c3)) {
		assert.Equal(t, http.StatusCreated, rec3.Code)
	}

}
