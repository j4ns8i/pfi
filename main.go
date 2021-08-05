package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// TODO setup CLI behavior
// with only one argument:
//   pfi should print out a palette generated from the argument which is
//   presumed to be an image
// add flags:
//   --input, -i to specify input template
//   --output, -o to specify rendered template destination
//   --config, -c to specify a config file which defines several templates to
//       render to destinations
//   --verbose, -v for verbose logging
//   --version because would it really be a CLI without it?

func main() {
	// trying out keeping these in main to prevent any other function besides
	// main from exiting the program except for panicking.
	usage := func() {
		fmt.Println(`usage: pfi INPUT OUTPUT IMAGE`)
		os.Exit(1)
	}
	exit := func(err error) {
		fmt.Println(err)
		os.Exit(1)
	}

	flag.Parse()
	inpath := flag.Arg(0)
	outpath := flag.Arg(1)
	imagepath := flag.Arg(2)
	switch {
	case inpath == "", outpath == "", imagepath == "":
		usage()
	}

	bytes, err := os.Open(imagepath)
	if err != nil {
		exit(err)
	}

	im, _, err := image.Decode(bytes)
	if err != nil {
		exit(err)
	}

	pal, err := generatePalette(im)
	if err != nil {
		exit(err)
	}

	if err = renderTemplate(inpath, outpath, pal); err != nil {
		exit(err)
	}

	printPalette(pal)
}
