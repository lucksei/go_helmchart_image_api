package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lucksei/go-chart-image-analyzer-api/internal/routes"
	"github.com/lucksei/go-chart-image-analyzer-api/internal/utils"
	"github.com/lucksei/go-chart-image-analyzer-api/middleware"
)

func main() {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	// NOTE: Initialize the result store BEFORE defining the methods, it will panic if not.
	// This is required to access the result store inside the endpoints
	resultStore := utils.NewResultStore()
	router.Use(middleware.ResultStore(resultStore))

	apiGroup := router.Group("/api")
	apiGroup.GET("/health", routes.Health)

	helmChartGroup := apiGroup.Group("/helm-chart")
	helmChartGroup.POST("", routes.HelmChartPost)
	helmChartGroup.GET("/:id", routes.HelmChartGet)

	// NOTE: Very important to initialize helm sdk settings before running API
	err := utils.InitHelmSettings()
	if err != nil {
		panic(err)
	}

	router.Run()

}
