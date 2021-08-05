package main

import (
	"fmt"

	"github.com/muesli/clusters"
)

type color struct {
	r float64
	g float64
	b float64
}

func (c color) String() string {
	return Hex(c)
}

// Return a color as hex, like #BADDAD.
func Hex(c color) string {
	return "#" + hexify(c.r, c.g, c.b)
}

// Return a color as hex without a leading # character, like BADDAD.
func AsStrippedHex(c color) string {
	return hexify(c.r, c.g, c.b)
}

type palette struct {
	Color0 color // Black
	Color1 color // Red
	Color2 color // Green
	Color3 color // Yellow
	Color4 color // Blue
	Color5 color // Magenta
	Color6 color // Cyan
	Color7 color // White
}

func (p palette) colors() [8]color {
	return [8]color{
		p.Color0,
		p.Color1,
		p.Color2,
		p.Color3,
		p.Color4,
		p.Color5,
		p.Color6,
		p.Color7,
	}
}

// Convert a [3]float64 representing normalized RGB values into three separate
// 8-bit RGB values.
func uint8ify(r, g, b float64) (uint8, uint8, uint8) {
	r255 := uint8(r * 0xff)
	g255 := uint8(g * 0xff)
	b255 := uint8(b * 0xff)
	return r255, g255, b255
}

// Convert a [3]float64 representing normalized RGB values into a hexadecimal
// string.
func hexify(r, g, b float64) string {
	r255, g255, b255 := uint8ify(r, g, b)
	return fmt.Sprintf("%X%X%X", r255, g255, b255)
}

// Normalize the given RGB values into float64 values between 0 and 1. The
// given values are between the range of [0,65535] according to the image/color
// package.
func observeRGB(r, g, b uint32) clusters.Observation {
	return clusters.Coordinates{
		float64(r) / float64(0xffff),
		float64(g) / float64(0xffff),
		float64(b) / float64(0xffff),
	}
}
