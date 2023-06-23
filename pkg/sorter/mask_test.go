package sorter

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskShouldCreateForValidMaskImage(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateMask(image)

	assert.Nil(t, err)
	assert.NotNil(t, mask)
}

func TestMaskShouldNotCreateForNilMaskImage(t *testing.T) {
	mask, err := CreateMask(nil)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldNotCreateForInvalidMaskImage(t *testing.T) {
	image := mockTestBlackAndRedImage()
	mask, err := CreateMask(image)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldCreateForValidMaskImageWithBounds(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateImageMask(image, image.Bounds(), 0)

	assert.Nil(t, err)
	assert.NotNil(t, mask)
}

func TestMaskShouldNotCreateForNilImageWithBounds(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateImageMask(nil, image.Bounds(), 0)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldNotCreateForInvalidBoundsWithBounds(t *testing.T) {
	img := mockTestBlackAndWhiteImage()
	mask, err := CreateImageMask(img, image.Rect(0, 0, 1, 1), 0)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldNotCreateForInvalidMaskImageWithBounds(t *testing.T) {
	image := mockTestBlackAndRedImage()
	mask, err := CreateImageMask(image, image.Bounds(), 0)

	assert.NotNil(t, err)
	assert.Nil(t, mask)
}

func TestMaskShouldCreateEmpty(t *testing.T) {
	mask := CreateEmptyMask()
	assert.NotNil(t, mask)
}

func TestMaskShouldTellIfIsMasked(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateMask(image)

	assert.Nil(t, err)
	assert.NotNil(t, mask)

	isMasked, err := mask.IsMasked(0, 0)
	assert.Nil(t, err)
	assert.False(t, isMasked)

	isMasked, err = mask.IsMasked(1, 0)
	assert.Nil(t, err)
	assert.True(t, isMasked)
}

func TestMaskShouldNotTellIfIsMaskedWhenLookupIsOutOfBounds(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateMask(image)

	assert.Nil(t, err)
	assert.NotNil(t, mask)

	_, err = mask.IsMasked(image.Bounds().Dx()+1, image.Bounds().Dy()+1)

	assert.NotNil(t, err)
}

func TestMaskShouldTellIfIsMaskedWhenTranslated(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateImageMask(image, image.Bounds(), 90)

	assert.Nil(t, err)
	assert.NotNil(t, mask)

	isMasked, err := mask.IsMasked(0, 0)
	assert.Nil(t, err)
	assert.True(t, isMasked)

	isMasked, err = mask.IsMasked(0, 1)
	assert.Nil(t, err)
	assert.False(t, isMasked)
}

func TestMaskShouldNotTellIfIsMaskedWhenLookupIsOutOfBoundsWhenTranslated(t *testing.T) {
	image := mockTestBlackAndWhiteImage()
	mask, err := CreateImageMask(image, image.Bounds(), 90)

	assert.Nil(t, err)
	assert.NotNil(t, mask)

	_, err = mask.IsMasked(image.Bounds().Dx()+1, image.Bounds().Dy()+1)

	assert.NotNil(t, err)
}

func TestMaskShouldTellIfIsMaskedForEmptyMask(t *testing.T) {
	mask := CreateEmptyMask()
	assert.NotNil(t, mask)

	// NOTE: Currently there is no index validation for empty mask, whatever value will return "False"
	isMasked, err := mask.IsMasked(0, 0)
	assert.Nil(t, err)
	assert.False(t, isMasked)
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
