package main

import (
	"github.com/betorvs/sonarqube-to-gitlab-webhook/config"
	"github.com/betorvs/sonarqube-to-gitlab-webhook/controller"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	e := echo.New()
	g := e.Group("/sonarqube-to-gitlab-webhook/v1")
	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	g.GET("/health", controller.CheckHealth)
	g.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	g.POST("/events", controller.ReceiveEvents)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
