package main

import (
	"flag"
	"fmt"

	"github.com/laghoule/png2pal/internal/pkg/img"

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
	flag.Parse()

	if *src == "" || *dst == "" || *gpl == "" {
		err := fmt.Errorf("png2pal -src <source file> -dst <destination file> -palette <GIMP palette file>")
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

	p := img.NewPalette()
	if err := p.Load(*gpl); err != nil {
		exitWithError(err)
	}

	imgDst := img.NewImage(*src, *dst, p)
	imgDst.Convert()

}

// exitWithError prints the error and exits with status code 1
func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
