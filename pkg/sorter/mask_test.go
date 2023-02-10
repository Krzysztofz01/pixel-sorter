package sorter

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Implement more mask color validation tests

func TestMaskShouldCreateForValidMaskImage(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateMask(image)

	assert.Nil(t, err)
	assert.NotNil(t, mask)
}

func TestMaskShouldNotCreateForInvalidMaskImage(t *testing.T) {
	image := mockTestBlackAndRedImage()
	mask, err := CreateMask(image)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldTellIfIsMasked(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateMask(image)
	assert.Nil(t, err)

	// NOTE: The test image consist of black and white vertical stripes, starting with the black color

	isMasked, err := mask.IsMasked(0, 0)
	assert.Nil(t, err)
	assert.False(t, isMasked)

	isMasked, err = mask.IsMasked(1, 0)
	assert.Nil(t, err)
	assert.True(t, isMasked)
}

// Create a test image which consists of valid mask colors (black and white)
func mockTestBlackAndWhiteImage() image.Image {
	const (
		height = 10
		width  = 10
	)

	p1 := image.Point{0, 0}
	p2 := image.Point{width, height}
	image := image.NewRGBA(image.Rectangle{p1, p2})

	cBlack := color.RGBA{0, 0, 0, 0xff}
	cWhite := color.RGBA{255, 255, 255, 0xff}

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			if xIndex%2 != 0 {
				image.Set(xIndex, yIndex, cBlack)
			} else {
				image.Set(xIndex, yIndex, cWhite)
			}
		}
	}

	return image
}

// Create a test image which consists a invalid mask color (red)
func mockTestBlackAndRedImage() image.Image {
	const (
		height = 10
		width  = 10
	)

	p1 := image.Point{0, 0}
	p2 := image.Point{width, height}
	image := image.NewRGBA(image.Rectangle{p1, p2})

	cBlack := color.RGBA{0, 0, 0, 0xff}
	cRed := color.RGBA{255, 0, 0, 0xff}

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			if xIndex%2 != 0 {
				image.Set(xIndex, yIndex, cBlack)
			} else {
				image.Set(xIndex, yIndex, cRed)
			}
		}
	}

	return image
}
