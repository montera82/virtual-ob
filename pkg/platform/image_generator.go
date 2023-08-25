package platform

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"time"
)

// GenerateRandomImageData creates an image with a random color and encodes it in PNG format.
// The generated image has a fixed width and height of 500 pixels.
// This function returns the bytes of the encoded image or an error if the encoding fails.
func GenerateRandomImageData() ([]byte, error) {
	rand.NewSource(time.Now().UnixNano())

	const width = 500
	const height = 500

	r := uint8(rand.Intn(256))
	g := uint8(rand.Intn(256))
	b := uint8(rand.Intn(256))
	randomColor := color.RGBA{r, g, b, 255}

	// Create the image with the random color.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, randomColor)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
