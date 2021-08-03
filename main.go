package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
)

func usage() {
	fmt.Println(`usage: pfi FILE`)
	os.Exit(1)
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	flag.Parse()
	input := flag.Arg(0)
	if input == "" {
		usage()
	}

	bytes, err := os.Open(input)
	if err != nil {
		exit(err)
	}

	im, _, err := image.Decode(bytes)
	if err != nil {
		exit(err)
	}

	bounds := im.Bounds()
	obs := make(clusters.Observations, 0, bounds.Dx()*bounds.Dy()/4)
	dsf := 16 // downsampling factor = 16 pixels per observation TODO do this smartly
	var r, g, b uint32
	for y := bounds.Min.Y; y < bounds.Max.Y; y += dsf {
		for x := bounds.Min.X; x < bounds.Max.X; x += dsf {
			r, g, b, _ = im.At(x, y).RGBA()
			obs = append(obs, observeRGB(r, g, b))
		}
	}

	km, err := kmeans.NewWithOptions(0.5, nil)
	if err != nil {
		exit(err)
	}

	// Partition the observations into 8 clusters, where each center represents
	// an approximation of a prominent color in the source.
	colors, err := km.Partition(obs, 8)
	if err != nil {
		exit(err)
	}

	// TODO build a proper output engine
	for i, c := range colors {
		fmt.Printf("color %v: #%s\n", i, hexify(c.Center))
	}
}

// Convert a [3]float64 representing normalized RGB values into a hexadecimal
// string.
func hexify(c clusters.Coordinates) string {
	r := uint8(c[0] * 0xff)
	g := uint8(c[1] * 0xff)
	b := uint8(c[2] * 0xff)
	return fmt.Sprintf("%X%X%X", r, g, b)
}

// Normalize the given rgb values into float64 values between 0 and 1. The
// given values are between the range of [0,65535] according to the image/color
// package.
func observeRGB(r, g, b uint32) clusters.Observation {
	return clusters.Coordinates{
		float64(r) / float64(0xffff),
		float64(g) / float64(0xffff),
		float64(b) / float64(0xffff),
	}
}
