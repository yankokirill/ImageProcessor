package main

import (
	"bytes"
	"encoding/base64"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"image/png"
	. "models"
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

func Process(task Task) {
	imageData, err := base64.StdEncoding.DecodeString(task.Payload.Image)
	if err != nil {
		commitTask(task.ID, "failed", "Invalid image data")
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		commitTask(task.ID, "failed", "Failed to decode image")
		return
	}

	switch task.Payload.Filter.Name {
	case "Grayscale":
		img = imaging.Grayscale(img)
	case "Blur":
		sigma, ok := task.GetFloatParameter("sigma")
		if !ok {
			commitTask(task.ID, "failed", "Invalid parameters")
			return
		}
		img = imaging.Blur(img, sigma)
	case "Sharpen":
		sigma, ok := task.GetFloatParameter("sigma")
		if !ok {
			commitTask(task.ID, "failed", "Invalid parameters")
			return
		}
		img = imaging.Sharpen(img, sigma)
	case "Negative":
		img = Negative(img)
	default:
		commitTask(task.ID, "failed", "Unknown filter")
		return
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		commitTask(task.ID, "failed", "Failed to encode image")
		return
	}
	encodedImage := base64.StdEncoding.EncodeToString(buf.Bytes())

	commitTask(task.ID, "ready", encodedImage)
}
