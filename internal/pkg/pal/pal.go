package pal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// Palette is a collection of colors
type Palette struct {
	Colors map[uint8]Color
}

// Color is a RGB color
type Color struct {
	R, G, B uint8
}

func NewPalette() *Palette {
	return &Palette{
		Colors: make(map[uint8]Color),
	}
}

func (p *Palette) Load(gplFile string) error {
	file, err := os.Open(gplFile)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewScanner(file)

	// Format: R G B comment
	re := regexp.MustCompile(`\d \d \d \w+`)

	for buf.Scan() {
		line := buf.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		parts := re.Split(line, -1)
		if len(parts) != 3 {
			continue
		}

		fmt.Println(parts)
	}

	return nil
}
