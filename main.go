package main

import (
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/controller"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Enable /metrics for prometheus
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	g := e.Group("/sonarqube-to-gitlab-webhook/v1")
	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	g.GET("/health", controller.CheckHealth)
	g.POST("/events", controller.ReceiveEvents)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
