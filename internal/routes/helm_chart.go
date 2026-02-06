package routes

import (
	"fmt"
	"net/http"
	"test/go_helm_chart_image_api/internal/models"
	"test/go_helm_chart_image_api/internal/utils"

	"github.com/gin-gonic/gin"
)

func HelmChartPost(c *gin.Context) {
	var jsonBody models.HelmChartRequest
	err := c.BindJSON(&jsonBody)
	if err != nil {
		c.Error(err)
	}

	helmChartPath := utils.HelmChartPath{
		RepoURL:   jsonBody.RepoURL,
		ChartPath: jsonBody.ChartPath,
	}
	helmChartId, err := helmChartPath.ToBase64Id()
	if err != nil {
		c.Error(err)
	}

	go func() {
		rendered, err := utils.RenderHelmTemplate(helmChartPath)
		if err != nil {
			panic(err)
		}

		images := utils.GetImagesFromRendered(rendered)
		for _, image := range images {
			imageInfo, err := utils.PullImageAndParseAPIInfo(image)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", imageInfo)
		}
	}()

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/helm-chart/%s", helmChartId))
	c.Status(http.StatusSeeOther)
}

func HelmChartGet(c *gin.Context) {
	id := c.Param("id")
	helmChartPath, err := utils.Base64StringToHelmChart(id)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusOK, models.HelmChartResponse{
		Status:    models.InProgress,
		RepoURL:   helmChartPath.RepoURL,
		ChartPath: helmChartPath.ChartPath,
		Images:    []models.ChartImage{},
	})
}
