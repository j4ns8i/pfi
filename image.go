package main

import (
	"errors"
	"image"

	"github.com/disintegration/imaging"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

const (
	// Maximum dimension of a downsampled image. Images are downsampled if
	// their greatest dimension is above this value, resulting in an image with
	// the same aspect ratio and the greatest dimension having this value.
	maxDownsampledDimension = 256
)

// Generate a color palette from an image.
func generatePalette(im image.Image) (palette, error) {
	bounds := im.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	switch {
	case w > maxDownsampledDimension && w <= h:
		im = imaging.Resize(im, 0, maxDownsampledDimension, imaging.Lanczos)
		bounds = im.Bounds()
	case h > maxDownsampledDimension && h <= w:
		im = imaging.Resize(im, maxDownsampledDimension, 0, imaging.Lanczos)
		bounds = im.Bounds()
	}

	obs := make(clusters.Observations, 0, bounds.Dx()*bounds.Dy())
	var r, g, b uint32
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ = im.At(x, y).RGBA()
			obs = append(obs, observeRGB(r, g, b))
		}
	}

	// Partition the observations into 8 clusters, where each center represents
	// an approximation of a prominent color in the source.
	km := kmeans.New()
	clusters, err := km.Partition(obs, 8)
	if err != nil {
		return palette{}, err
	}

	if len(clusters) != 8 {
		panic(errors.New("Expected 8 clusters from data"))
	}

	// TODO sort colors into the most appropriate mapping based on lowest total
	// std dev, or something?
	pal := palette{
		Color0: colorFromCluster(clusters[0]),
		Color1: colorFromCluster(clusters[1]),
		Color2: colorFromCluster(clusters[2]),
		Color3: colorFromCluster(clusters[3]),
		Color4: colorFromCluster(clusters[4]),
		Color5: colorFromCluster(clusters[5]),
		Color6: colorFromCluster(clusters[6]),
		Color7: colorFromCluster(clusters[7]),
	}
	return pal, nil
}

// Convert a k-means cluster of [3]float64 values into three separate
// normalized RGB values between 0 and 1, representing the cluster's "center",
// i.e. the color which best represents that cluster.
func colorFromCluster(cluster clusters.Cluster) color {
	center := cluster.Center
	if len(center) < 3 {
		// TODO add testing to catch this scenario since it shouldn't depend on
		// user input
		panic(errors.New("expected 3 values as RGB data"))
	}
	return color{r: center[0], g: center[1], b: center[2]}
}
