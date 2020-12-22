package controller

import (
	"net/http"

	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/domain"
	"github.com/labstack/echo/v4"
)

//GetInfo of the application like version
func GetInfo(c echo.Context) error {
	info := domain.Info{Version: config.Version}
	return c.JSON(http.StatusOK, info)
}
