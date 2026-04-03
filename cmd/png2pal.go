package main

import (
	"flag"
	"fmt"

	"github.com/laghoule/png2pal/internal/pkg/pal"

	"image/color"
	"image/png"
	"os"
)

var (
	version   = "unknown"
	gitCommit = "unknown"
)

func main() {
	src := flag.String("src", "", "source file")
	dst := flag.String("dst", "x", "destination file")
	gpl := flag.String("palette", "mia.gpl", "GIMP palette file")
	idx := flag.Int("index", 0, "index of the palette to use for transparency")
	flag.Parse()

	if *src == "" || *dst == "" || *gpl == "" {
		err := fmt.Errorf("png2pal -src <source file> -dst <destination file> -palette <GIMP palette file>")
		exitWithError(err)
	}

	if *idx < 0 || *idx > 255 {
		err := fmt.Errorf("png2pal: index must be between 0 and 255")
		exitWithError(err)
	}

	fileSrc, err := os.Open(*src)
	if err != nil {
		exitWithError(err)
	}
	defer fileSrc.Close()

	imgSrc, err := png.Decode(fileSrc)
	if err != nil {
		exitWithError(err)
	}

	// NRGBAModel is the color model for RGBA images.
	if imgSrc.ColorModel() != color.NRGBAModel {
		err := fmt.Errorf("png2pal: source image must be in RGBA format")
		exitWithError(err)
	}

	p := pal.NewPalette()
	if err := p.Load(*gpl); err != nil {
		exitWithError(err)
	}

	for i := 0; i < len(p.Colors); i++ {
		if i == *idx {
			fmt.Printf("Transparent index %d: %v\n", i, p.Colors[uint8(i)])
		}
	}

	testColor := pal.Color{R: 25,G: 50,B: 100}
	fmt.Printf("Test color: %v\n", testColor)
	closestIndex := p.FindClosestColorIndex(testColor)
	fmt.Printf("Closest color: %v at index %d\n", p.Colors[closestIndex], closestIndex)

}

// exitWithError prints the error and exits with status code 1
func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
