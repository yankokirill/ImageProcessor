package filter

import (
	"image"
	"image/color"
)

func Negative(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	negativeImg := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			negativeImg.Set(bounds.Max.X-(x-bounds.Min.X+1), y, color.NRGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255,
			})
		}
	}
	return negativeImg
}
