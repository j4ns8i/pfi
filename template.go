package main

// TODO build a proper output engine

import (
	"fmt"
	"html/template"
	"os"
)

// Render a template at inpath to the given outpath using palette pal.
func renderTemplate(inpath, outpath string, pal palette) error {
	// TODO add functions like Hex, StrippedHex, R, G, B via template.FuncMap
	t, err := template.ParseFiles(inpath)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(outpath, os.O_RDWR|os.O_CREATE, 0666) // mode is applied *before* umask
	if err != nil {
		return err
	}
	if err = t.Execute(f, pal); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	return nil
}

// Print out the clusters' centers, which represented normalized RGB values, as
// colors on the terminal.
func printPalette(p palette) {
	height := 3
	width := 5
	var r, g, b uint8
	for y := 0; y < height; y++ {
		for _, c := range p.colors() {
			for x := 0; x < width; x++ {
				r, g, b = uint8ify(c.r, c.g, c.b)
				fmt.Printf("\x1b[0;48;2;%d;%d;%dm ", r, g, b)
			}
			fmt.Printf("\x1b[0m ")
		}
		fmt.Printf("\x1b[0m\n")
	}
}
