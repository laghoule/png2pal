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
	dst := flag.String("dst", "", "destination file")
	pal := flag.String("pal", "", "GIMP palette file")
	flag.Parse()

	fmt.Printf("png2pal version: %s, git commit: (%s)\n\n", version, gitCommit)

	if *src == "" || *dst == "" || *pal == "" {
		err := fmt.Errorf("png2pal -src <source file> -dst <destination file> -pal <GIMP palette file>")
		exitWithError(err)
	}

	fmt.Printf("Converting %s to %s\n", *src, *dst)
	img, err := img.NewImage(*src, *dst, *pal)
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
