package pal

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// Palette is a collection of colors
type Palette struct {
	Colors map[uint8]Color
}

// Color is a RGB color
type Color struct {
	R, G, B uint8
}

// NewPalette creates a new Palette
func NewPalette() *Palette {
	return &Palette{
		Colors: make(map[uint8]Color, 255),
	}
}

// Load loads a GIMP palette file into the Palette
func (p *Palette) Load(gplFile string) error {
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
			return fmt.Errorf("incolid color at line %d: %w", j, err)
		}
		if color == nil {
			continue
		}

		p.Colors[i] = *color
		i++
	}

	return nil
}

// extractRGB extracts a RGB color from a line of text
func extractRGB(re *regexp.Regexp, line string) (*Color, error) {
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return nil, nil
	}

	r, g, b, err := convertRGBMatchesToInt(matches[1:])
	if err != nil {
		return nil, err
	}

	return &Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
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

// FindClosestColorIndex finds the index of the closest color in the palette to the target color
// Euclidian distance: D = √ [ (R1 - R2)² + (G1 - G2)² + (B1 - B2)² ]
func (p *Palette) FindClosestColorIndex(target Color) uint8 {
	var closestIndex uint8
	minDistSq := -1 // initialize to invalid value (for first comparison)

	for idx, palColor := range p.Colors {
		dr := int(palColor.R) - int(target.R)
		dg := int(palColor.G) - int(target.G)
		db := int(palColor.B) - int(target.B)

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
