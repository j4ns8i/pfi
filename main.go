package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/disintegration/imaging"
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
	w := bounds.Dx()
	h := bounds.Dy()
	dsm := 256 // the max dimension after downsampling
	switch {
	case w > dsm && w <= h:
		im = imaging.Resize(im, 0, dsm, imaging.Lanczos)
		bounds = im.Bounds()
	case h > dsm && h <= w:
		im = imaging.Resize(im, dsm, 0, imaging.Lanczos)
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
	colors, err := km.Partition(obs, 8)
	if err != nil {
		exit(err)
	}

	printPalette(colors)
}

// Print out the clusters' centers, which represented normalized rgb values, as
// colors on the terminal.
//
// TODO build a proper output engine
func printPalette(c clusters.Clusters) {
	height := 2
	width := 5
	for y := 0; y < height; y++ {
		for _, color := range c {
			r, g, b := rgb255(color.Center)
			for x := 0; x < width; x++ {
				fmt.Printf("\x1b[0;48;2;%d;%d;%dm ", r, g, b)
			}
			fmt.Printf("\x1b[0m ")
		}
		fmt.Printf("\x1b[0m\n")
	}
}

// Convert a [3]float64 representing normalized rgb values into three separate 8-bit rgb values
func rgb255(c clusters.Coordinates) (uint8, uint8, uint8) {
	r := uint8(c[0] * 0xff)
	g := uint8(c[1] * 0xff)
	b := uint8(c[2] * 0xff)
	return r, g, b
}

// Convert a [3]float64 representing normalized rgb values into a hexadecimal
// string.
func hexify(c clusters.Coordinates) string {
	r, g, b := rgb255(c)
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
