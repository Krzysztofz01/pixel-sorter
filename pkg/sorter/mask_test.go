package sorter

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskShouldCreateForValidRgbaMaskImage(t *testing.T) {
	mask, err := CreateMaskFromRgba(mockMaskTestImageRgba(whiteRgba, blackRgba))

	assert.Nil(t, err)
	assert.NotNil(t, mask)
}

func TestMaskShouldNotCreateForInvalidRgbaMaskImage(t *testing.T) {
	mask, err := CreateMaskFromRgba(mockMaskTestImageRgba(whiteRgba, redRgba))

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldNotCreateForNilRgbaMaskImage(t *testing.T) {
	mask, err := CreateMaskFromRgba(nil)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldCreateForValidNrgbaMaskImage(t *testing.T) {
	mask, err := CreateMaskFromNrgba(mockMaskTestImageNrgba(whiteNrgba, blackNrgba))

	assert.Nil(t, err)
	assert.NotNil(t, mask)
}

func TestMaskShouldNotCreateForInvalidNrgbaMaskImage(t *testing.T) {
	mask, err := CreateMaskFromNrgba(mockMaskTestImageNrgba(whiteNrgba, redNrgba))

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldNotCreateForNilNrgbaMaskImage(t *testing.T) {
	mask, err := CreateMaskFromNrgba(nil)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldCreateEmpty(t *testing.T) {
	mask := CreateEmptyMask()
	assert.NotNil(t, mask)
}

func TestMaskShouldCorrectlyReturnAt(t *testing.T) {
	image := mockMaskTestImageNrgba(whiteNrgba, blackNrgba)
	mask, err := CreateMaskFromNrgba(image)

	assert.Nil(t, err)
	assert.NotNil(t, mask)

	at, err := mask.At(0, 0)
	assert.Nil(t, err)
	assert.Equal(t, uint8(0x00), at)

	at, err = mask.At(1, 0)
	assert.Nil(t, err)
	assert.Equal(t, uint8(0xff), at)

	atb, err := mask.AtB(0, 0)
	assert.Nil(t, err)
	assert.Equal(t, true, atb)

	atb, err = mask.AtB(1, 0)
	assert.Nil(t, err)
	assert.Equal(t, false, atb)

	at, err = mask.AtByIndex(0)
	assert.Nil(t, err)
	assert.Equal(t, uint8(0x00), at)

	at, err = mask.AtByIndex(1)
	assert.Nil(t, err)
	assert.Equal(t, uint8(0xff), at)

	atb, err = mask.AtByIndexB(0)
	assert.Nil(t, err)
	assert.Equal(t, true, atb)

	atb, err = mask.AtByIndexB(1)
	assert.Nil(t, err)
	assert.Equal(t, false, atb)
}

func TestMaskShouldReturnErrorOnInvalidAt(t *testing.T) {
	image := mockMaskTestImageNrgba(whiteNrgba, blackNrgba)
	mask, err := CreateMaskFromNrgba(image)

	assert.Nil(t, err)
	assert.NotNil(t, mask)

	_, err = mask.At(-1, -1)
	assert.NotNil(t, err)

	_, err = mask.At(testMaskImageWidth, testMaskImageHeight)
	assert.NotNil(t, err)

	_, err = mask.AtB(-1, -1)
	assert.NotNil(t, err)

	_, err = mask.AtB(testMaskImageWidth, testMaskImageHeight)
	assert.NotNil(t, err)

	_, err = mask.AtByIndex(-1)
	assert.NotNil(t, err)

	_, err = mask.AtByIndex(testMaskImageWidth * testMaskImageHeight)
	assert.NotNil(t, err)

	_, err = mask.AtByIndexB(-1)
	assert.NotNil(t, err)

	_, err = mask.AtByIndexB(testMaskImageWidth * testMaskImageHeight)
	assert.NotNil(t, err)
}

const (
	testMaskImageHeight int = 10
	testMaskImageWidth  int = 10
)

var (
	whiteRgba  color.RGBA  = color.RGBA{0xff, 0xff, 0xff, 0xff}
	blackRgba  color.RGBA  = color.RGBA{0x00, 0x00, 0x00, 0xff}
	redRgba    color.RGBA  = color.RGBA{0x0ff, 0x00, 0x00, 0x0ff}
	whiteNrgba color.NRGBA = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	blackNrgba color.NRGBA = color.NRGBA{0x00, 0x00, 0x00, 0xff}
	redNrgba   color.NRGBA = color.NRGBA{0x0ff, 0x00, 0x00, 0x0ff}
)

// Create a test image which consists of valid mask colors (black and white)
func mockMaskTestImageRgba(a, b color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, testMaskImageWidth, testMaskImageHeight))

	for xIndex := 0; xIndex < testMaskImageWidth; xIndex += 1 {
		for yIndex := 0; yIndex < testMaskImageHeight; yIndex += 1 {
			if xIndex%2 != 0 {
				img.SetRGBA(xIndex, yIndex, a)
			} else {
				img.SetRGBA(xIndex, yIndex, b)
			}
		}
	}

	return img
}

func mockMaskTestImageNrgba(a, b color.NRGBA) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, testMaskImageWidth, testMaskImageHeight))

	for xIndex := 0; xIndex < testMaskImageWidth; xIndex += 1 {
		for yIndex := 0; yIndex < testMaskImageHeight; yIndex += 1 {
			if xIndex%2 != 0 {
				img.SetNRGBA(xIndex, yIndex, a)
			} else {
				img.SetNRGBA(xIndex, yIndex, b)
			}
		}
	}

	return img
}
