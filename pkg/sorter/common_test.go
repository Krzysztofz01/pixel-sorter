package sorter

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateMaxIntervalLengthShouldCalculateCorrectIntervalMaxLengthForGivenOptions(t *testing.T) {
	cases := []struct {
		length       int
		randomFactor int
		expected     int
	}{
		{0, 0, 0},
	}

	for _, c := range cases {
		actual := calculateMaxIntervalLength(c.length, c.randomFactor)

		assert.Equal(t, c.expected, actual)
	}
}

func TestDrawBufferToImageShouldCorrectlyAppendBufferToImage(t *testing.T) {
	bounds := image.Rect(0, 0, 8, 1)
	actualImg := image.NewRGBA(bounds)
	expectedImg := image.NewRGBA(bounds)

	for x := 0; x < bounds.Dx(); x += 1 {
		if x >= 2 && x <= 5 {
			expectedImg.Set(x, 0, color.White)
		} else {
			expectedImg.Set(x, 0, color.Black)
		}

		actualImg.Set(x, 0, color.Black)
	}

	buffer := []color.RGBA{
		{0xff, 0xff, 0xff, 0xff},
		{0xff, 0xff, 0xff, 0xff},
		{0xff, 0xff, 0xff, 0xff},
		{0xff, 0xff, 0xff, 0xff},
	}

	drawBufferIntoImage(actualImg, buffer, 5*4, 1*4)

	assert.Equal(t, expectedImg.Pix, actualImg.Pix)
}
