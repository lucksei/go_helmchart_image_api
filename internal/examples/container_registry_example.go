package examples

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func ContainerRegistryCustomExample() {
	fmt.Println("Parsing image reference...")
	ref, err := name.ParseReference("docker.io/kooldev/pause:latest")
	// ref, err := name.ParseReference("docker.io/registry:2")
	if err != nil {
		panic(err)
	}
	fmt.Println("Pulling image...")
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		panic(err)
	}
	if s, err := img.Size(); err != nil {
		panic(err)
	} else {
		fmt.Printf("Size (from image): %d\n", s)
	}

	manifest, err := img.Manifest()
	fmt.Printf("Size (from manifest.Config.Size): %d\n", manifest.Config.Size)
	for i, layer := range manifest.Layers {
		digest := layer.Digest
		size := layer.Size
		fmt.Printf("%d: %s (%d) \n", i, digest.String(), size)
	}
}
