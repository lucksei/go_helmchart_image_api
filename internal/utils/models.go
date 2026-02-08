package utils

type ImageAnalysis struct {
	Name           string `json:"name"`
	Size           int64  `json:"size"`
	NumberOfLayers int    `json:"no_of_layers"`
}

type HelmChartSource struct {
	RepoURL  string `json:"repo_url"`
	ChartRef string `json:"chart_ref" binding:"required"`
}

type HelmChartAnalysis struct {
	RepoURL  string          `json:"repo_url"`
	ChartRef string          `json:"chart_path"`
	Images   []ImageAnalysis `json:"images"`
}
