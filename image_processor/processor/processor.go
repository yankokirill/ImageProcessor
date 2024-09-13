package processor

import (
	"bytes"
	"encoding/base64"
	"github.com/disintegration/imaging"
	. "hw/image_processor/my_filters"
	. "hw/models"
	"image"
	"image/png"
)

func Process(task Task) (status, result string) {
	imageData, err := base64.StdEncoding.DecodeString(task.Payload.Image)
	if err != nil {
		return "failed", "Invalid image data"
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "failed", "Failed to decode image"
	}

	switch task.Payload.Filter.Name {
	case "Grayscale":
		img = imaging.Grayscale(img)
	case "Blur":
		sigma, ok := task.GetFloatParameter("sigma")
		if !ok {
			return "failed", "Invalid parameters"
		}
		img = imaging.Blur(img, sigma)
	case "Sharpen":
		sigma, ok := task.GetFloatParameter("sigma")
		if !ok {
			return "failed", "Invalid parameters"
		}
		img = imaging.Sharpen(img, sigma)
	case "Negative":
		img = Negative(img)
	default:
		return "failed", "Unknown filter"
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "failed", "Failed to encode image"
	}
	result = base64.StdEncoding.EncodeToString(buf.Bytes())
	return "ready", result
}
