package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/usecase"
	"github.com/labstack/echo"
)

//ReceiveEvents func
func ReceiveEvents(c echo.Context) (err error) {
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	bodyString := string(bodyBytes)
	event := new(domain.Events)
	if err = c.Bind(event); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	// Headers
	sonarWebhook := c.Request().Header.Get("X-Sonar-Webhook-Hmac-Sha256")
	verifier := usecase.ValidateWebhook(sonarWebhook, bodyString)
	if verifier != true {
		return c.JSON(http.StatusForbidden, nil)
	}
	go usecase.GitlabCommit(event.Project.Name, event.Revision, event.Project.URL, event.Status)
	go fmt.Printf("[INFO]: Project Name %s, Project URL: %s, Status: %s, Revision: %s", event.Project.Name, event.Project.URL, event.Status, event.Revision)
	return c.JSON(http.StatusCreated, "OK")
}
