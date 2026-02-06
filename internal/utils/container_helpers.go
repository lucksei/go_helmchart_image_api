package utils

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

type ImageInfo struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	LayerNumber int    `json:"layer_number"`
}

func PullImageAndParseAPIInfo(image string) (ImageInfo, error) {
	imageInfo := ImageInfo{}

	ref, err := name.ParseReference(image)
	if err != nil {
		return ImageInfo{}, err
	}
	imageInfo.Name = ref.Name()

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return ImageInfo{}, err
	}
	manifest, err := img.Manifest()
	if err != nil {
		return ImageInfo{}, err
	}

	imageInfo.LayerNumber = len(manifest.Layers)

	for _, layer := range manifest.Layers {
		imageInfo.Size += layer.Size
	}

	return imageInfo, err
}
