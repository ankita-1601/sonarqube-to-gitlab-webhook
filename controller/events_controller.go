package controller

import (
	"bytes"
	"io/ioutil"
	"log"
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
	bodyString := string(bodyBytes)
	// Headers Validation
	sonarWebhook := c.Request().Header.Get("X-Sonar-Webhook-Hmac-Sha256")
	verifier := usecase.ValidateWebhook(sonarWebhook, bodyString)
	if !verifier {
		return c.JSON(http.StatusForbidden, nil)
	}
	// Use the content
	event := new(domain.Events)
	if err = c.Bind(event); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	go log.Printf("[INFO]: Project Name %s, Project URL: %s, Status: %s, Revision: %s", event.Project.Name, event.Project.URL, event.Status, event.Revision)
	go usecase.GitlabCommit(event.Project.Name, event.Revision, event.Project.URL, event.Status)

	return c.JSON(http.StatusCreated, "OK")
}
