package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/laghoule/png2pal/internal/pkg/img"
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

	p := img.NewPalette()
	if err := p.Load(*gpl); err != nil {
		exitWithError(err)
	}

	img, err := img.NewImage(*src, *dst, *gpl)
	if err != nil {
		exitWithError(err)
	}
	
	if err := img.Convert(); err != nil {
		exitWithError(err)
	}

}

// exitWithError prints the error and exits with status code 1
func exitWithError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
