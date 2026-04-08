package img

import (
	"fmt"
	"image"
	gocolor "image/color"
	"image/png"
	"os"
)

type img struct {
	src string
	dst string
	pal *palette
}

// NewImage creates a new img instance
func NewImage(src, dst, pal string) (*img, error) {
	p := NewPalette()
	if err := p.Load(pal); err != nil {
		return nil, fmt.Errorf("png2pal: failed to load palette: %v", err)
	}

	return &img{
		src: src,
		dst: dst,
		pal: p,
	}, nil
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

	// NRGBAModel is the color model for RGBA images.
	if srcImage.ColorModel() != gocolor.NRGBAModel {
		return fmt.Errorf("png2pal: source image must be in RGBA format")
	}

	newRect := srcImage.Bounds()
	destImage := image.NewPaletted(newRect, i.pal.ToColorPaletted())
	c := gocolor.RGBA{}

	for y := range newRect.Max.Y {
		for x := range newRect.Max.X {
			r, g, b, a := srcImage.At(x, y).RGBA()

			// if alpha is 0, set to transparent
			if a == 0 {
				destImage.SetColorIndex(x, y, 0)
				continue
			}

			// convert to 8 bit color
			c = gocolor.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
			}

			// find the nearest color in the palette
			nearestColorIndex := i.pal.FindClosestColorIndex(color{
				r: c.R,
				g: c.G,
				b: c.B,
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
