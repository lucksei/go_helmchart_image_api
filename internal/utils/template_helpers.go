package utils

import (
	"strings"

	"github.com/goccy/go-yaml"
)

func ContainersSpecSearch(m map[string]any, res *[]any) {
	for k, v := range m {
		if k == "containers" {
			*res = append(*res, v.([]any))
		}
		if m, ok := v.(map[string]any); ok {
			ContainersSpecSearch(m, res)
		}
	}
}

func GetImagesFromRendered(r map[string]string) []string {
	var images = []string{}
	for key, value := range r {
		if strings.Contains(key, ".yaml") {
			var template map[string]any
			yaml.Unmarshal([]byte(value), &template)

			var containerSpecList []any
			ContainersSpecSearch(template, &containerSpecList)
			for _, containers := range containerSpecList {
				for _, container := range containers.([]any) {
					images = append(images, container.(map[string]any)["image"].(string))
				}
			}
		}
	}
	return images
}
