package main

import (
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/controller"
	_ "github.com/betorvs/sonarqube-to-gitlab-webhook/gateway/customlog"
	_ "github.com/betorvs/sonarqube-to-gitlab-webhook/gateway/gitlabclient"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	controller.MapRoutes(e)

	e.Logger.Fatal(e.Start(":" + config.Values.Port))
}
