package routes

import (
	"fmt"
	"net/http"
	"test/go_helm_chart_image_api/internal/models"
	"test/go_helm_chart_image_api/internal/utils"

	"github.com/gin-gonic/gin"
)

func HelmChartPost(c *gin.Context) {
	// Validating the request
	var jsonBody models.HelmChartRequest
	err := c.BindJSON(&jsonBody)
	if err != nil {
		c.Error(err)
	}

	// Create an ID for the specific helm chart (TODO: Improve this for different import methods???)
	helmChartSource := utils.HelmChartSource{
		RepoURL:  jsonBody.RepoURL,
		ChartRef: jsonBody.ChartPath,
	}
	helmChartId, err := helmChartSource.ToBase64Id()
	if err != nil {
		c.Error(err)
	}

	// Loading the result store
	rs, ok := c.MustGet("result_store").(*utils.ResultStore)
	if !ok {
		c.Error(err)
	}

	// If the helm chart is being processed, accept (202)
	_, status := rs.Get(helmChartId)
	if status == utils.StatusInProgress {
		c.Writer.Header().Set("Location", fmt.Sprintf("/api/helm-chart/%s", helmChartId))
		c.Status(http.StatusAccepted)
		return
	}
	// If the helm chart is already in the store, redirect (303)
	if status == utils.StatusSuccess {
		c.Writer.Header().Set("Location", fmt.Sprintf("/api/helm-chart/%s", helmChartId))
		c.Status(http.StatusSeeOther)
		return
	}

	// TODO: Add way of knowing if its already processing to prevent parallel requests running the same thing tiwce
	// Runs in the background, processes the helm chart + images if they are not already inside the store
	rs.SetPending(helmChartId)
	go func() {
		fmt.Printf("Processing helm chart %s\n", helmChartSource.ChartRef)
		rendered, err := utils.RenderHelmTemplate(helmChartSource)
		if err != nil {
			panic(err)
		}

		images := utils.GetImagesFromRendered(rendered)
		fmt.Printf("Found %d images\n", len(images))
		fmt.Printf("Processing images\n")
		result := utils.HelmChartAnalysis{
			RepoURL:  helmChartSource.RepoURL,
			ChartRef: helmChartSource.ChartRef,
		}
		imagesAnalysis := []utils.ImageAnalysis{}
		for i, image := range images {
			imageAnalysis, err := utils.PullImageAndParseAPIInfo(image)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Image %d: %s\n", i, imageAnalysis.Name)
			imagesAnalysis = append(imagesAnalysis, imageAnalysis)
		}
		result.Images = imagesAnalysis
		fmt.Printf("%v\n", result)

		rs.Put(helmChartId, result)
		fmt.Printf("Done processing helm chart %s\n", helmChartSource.ChartRef)
	}()

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/helm-chart/%s", helmChartId))
	c.Status(http.StatusAccepted)
}

func HelmChartGet(c *gin.Context) {
	id := c.Param("id")
	// helmChartSource, err := utils.Base64StringToHelmChart(id)
	// if err != nil {
	// 	c.Error(err)
	// }

	// Loading the result store
	rs, ok := c.MustGet("result_store").(*utils.ResultStore)
	if !ok {
		c.Error(fmt.Errorf("Failed to retrieve ResultStore"))
	}
	result, status := rs.Get(id)
	if status == utils.StatusInProgress {
		return
	}
	if status == utils.StatusNotFound {
		c.Status(http.StatusNotFound)
		return
	}

	fmt.Printf("TODO make this into a response later:\n%v\n", result)
	c.JSON(http.StatusOK, result)
	// TODO: Ugly conversion, fix or move to helper inside the models
	// result, ok := rs.Get(id)
	// if !ok {
	// 	c.JSON(http.StatusOK, models.HelmChartResponse{
	// 		Status:    models.InProgress,
	// 		RepoURL:   helmChartPath.RepoURL,
	// 		ChartPath: helmChartPath.ChartRef,
	// 		Images:    []models.ChartImage{},
	// 	})
	// 	return
	// }
	// response := models.HelmChartResponse{
	// 	Status:    models.Success,
	// 	RepoURL:   result.RepoURL,
	// 	ChartPath: result.ChartPath,
	// }
	// images := []models.ChartImage{}
	// for _, img := range result.Images {
	// 	images = append(images, models.ChartImage{
	// 		Name:        img.Name,
	// 		Size:        img.Size,
	// 		LayerNumber: img.LayerNumber,
	// 	})
	// }
	// response.Images = images
	// c.JSON(http.StatusOK, response)
}
