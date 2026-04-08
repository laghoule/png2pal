package img

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type img struct {
	src string
	dst string
	pal Palette
}

// NewImage creates a new img instance
func NewImage(src, dst string, pal *Palette) *img {
	return &img{
		src: src,
		dst: dst,
		pal: *pal,
	}
}

func (i *img) Convert() error {
	dstFile, err := os.Create(i.dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(i.src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcImage, err := png.Decode(srcFile)
	if err != nil {
		return err
	}

	newRect := srcImage.Bounds()
	destImage := image.NewPaletted(newRect, i.pal.ToColorPaletted())
	c := color.RGBA{}

	for y := range newRect.Max.Y {
		for x := range newRect.Max.X {
			r, g, b, a := srcImage.At(x, y).RGBA()

			// if alpha is 0, set to transparent
			if a == 0 {
				destImage.SetColorIndex(x, y, 0)
				continue
			}

			// convert to 8 bit color
			c = color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
			}

			// find the nearest color in the palette
			nearestColorIndex := i.pal.FindClosestColorIndex(Color{
				R: c.R,
				G: c.G,
				B: c.B,
			})

			destImage.SetColorIndex(x, y, nearestColorIndex)
		}

	}

	err = png.Encode(dstFile, destImage)
	if err != nil {
		return err
	}

	return nil
}
