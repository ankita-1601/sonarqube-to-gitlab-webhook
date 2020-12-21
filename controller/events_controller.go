package controller

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/usecase"
	"github.com/labstack/echo/v4"
)

//ReceiveEvents func
func ReceiveEvents(c echo.Context) (err error) {
	if config.Values.SonarqubeSecret == "Absent" {
		return c.JSON(http.StatusNotImplemented, nil)
	}
	var bodyBytes []byte
	if c.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyString := string(bodyBytes)
	// Headers Validation
	sonarWebhook := c.Request().Header.Get("X-Sonar-Webhook-Hmac-Sha256")
	verifier := usecase.ValidateWebhook(sonarWebhook, bodyString, config.Values.SonarqubeSecret)
	if !verifier {
		return c.JSON(http.StatusForbidden, nil)
	}
	// Use the content
	event := new(domain.Events)
	if err = c.Bind(event); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	logger := config.GetLogger
	defer logger().Sync()
	logger().Infof("Project Name %s, Project URL %s, Status %s, Revision %s", event.Project.Name, event.Project.URL, event.Status, event.Revision)
	err = usecase.GitlabCommit(event)
	if err != nil {
		logger().Info("cannot post a commit to gitlab")
	}

	return c.JSON(http.StatusCreated, "OK")
}
