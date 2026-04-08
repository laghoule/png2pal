package img

import (
	"bufio"
	"fmt"
	gocolor "image/color"
	"os"
	"regexp"
	"strconv"
)

const (
	transparantAlpha = 0
	opaqueAlpha      = 255
	VGAColors        = 256
	lastIndexedColor = VGAColors - 1
)

// Palette is a collection of colors
type palette struct {
	colors map[uint8]color
}

// Color is a RGB color
type color struct {
	r, g, b uint8
}

// NewPalette creates a new Palette
func NewPalette() *palette {
	return &palette{
		colors: make(map[uint8]color, 255),
	}
}

// Load loads a GIMP palette file into the Palette
func (p *palette) Load(gplFile string) error {
	file, err := os.Open(gplFile)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewScanner(file)

	// Format: R G B comment
	re := regexp.MustCompile(`^(\d+)\s+(\d+)\s+(\d+)\s+(.*)$`)

	var i uint8
	var j int
	for buf.Scan() {
		j++
		line := buf.Text()
		if line == "" || line[0] == '#' {
			continue
		}

		color, err := extractRGB(re, line)
		if err != nil {
			return fmt.Errorf("invalid color at line %d: %w", j, err)
		}
		if color == nil {
			continue
		}

		p.colors[i] = *color
		i++
	}

	if len(p.colors) < VGAColors {
		return fmt.Errorf("invalid number of colors: %d, should be %d", len(p.colors), VGAColors)
	}

	return nil
}

// extractRGB extracts a RGB color from a line of text
func extractRGB(re *regexp.Regexp, line string) (*color, error) {
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	r, g, b, err := convertRGBMatchesToInt(matches[1:])
	if err != nil {
		return nil, err
	}

	return &color{
		r: uint8(r),
		g: uint8(g),
		b: uint8(b),
	}, nil
}

// convertRGBMatchesToInt converts a slice of RGB matches to int values
func convertRGBMatchesToInt(matches []string) (int, int, int, error) {
	if len(matches) != 4 {
		return 0, 0, 0, fmt.Errorf("incolid number of matches: %d", len(matches))
	}

	var col [3]int
	var err error

	for i := range col {
		col[i], err = strconv.Atoi(matches[i])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid value at index %d: %v", i, err)
		}

		if col[i] < 0 || col[i] > 255 {
			return 0, 0, 0, fmt.Errorf("invalid value at index %d: %d", i, col[i])
		}
	}

	return col[0], col[1], col[2], nil
}

// findClosestColorIndex finds the index of the closest color in the palette to the target color
// Euclidian distance: D = √ [ (R1 - R2)² + (G1 - G2)² + (B1 - B2)² ]
func (p *palette) findClosestColorIndex(target color) uint8 {
	var closestIndex uint8
	minDistSq := -1 // initialize to invalid value (for first comparison)

	for idx, palColor := range p.colors {
		if idx == 0 {
			continue // skip the transparent color
		}

		dr := int(palColor.r) - int(target.r)
		dg := int(palColor.g) - int(target.g)
		db := int(palColor.b) - int(target.b)

		drSq := dr * dr
		dgSq := dg * dg
		dbSq := db * db

		// No need to calculate the square root, since we are comparing distances
		currentDist := int(drSq + dgSq + dbSq)

		if minDistSq == -1 || currentDist < minDistSq {
			minDistSq = currentDist
			closestIndex = uint8(idx)
		}
	}

	return closestIndex
}

// toColorPaletted converts the Palette to a color.Palette
func (p *palette) toColorPaletted() gocolor.Palette {
	colPal := make(gocolor.Palette, VGAColors)

	// Set the first color to transparent
	colPal[0] = gocolor.RGBA{
		R: p.colors[uint8(0)].r,
		G: p.colors[uint8(0)].g,
		B: p.colors[uint8(0)].b,
		A: transparantAlpha,
	}

	for i := 1; i <= lastIndexedColor; i++ {
		if len(p.colors) >= i {
			colPal[i] = gocolor.RGBA{
				R: p.colors[uint8(i)].r,
				G: p.colors[uint8(i)].g,
				B: p.colors[uint8(i)].b,
				A: opaqueAlpha,
			}
			continue
		}

		// If the palette is smaller than the expected size, fill the rest with the first color
		colPal[i] = colPal[0]
	}

	return colPal
}
