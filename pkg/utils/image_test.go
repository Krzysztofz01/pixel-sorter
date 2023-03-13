package utils

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetColumnShouldReturnImageColumn(t *testing.T) {
	image := mockTestGradientImage()

	lookupMap := make(map[color.Color]bool)
	expectedLength := 1

	column, err := GetImageColumn(image, 0)
	assert.Nil(t, err)

	for _, color := range column {
		lookupMap[color] = true
	}

	acutalLength := len(lookupMap)
	assert.Equal(t, expectedLength, acutalLength)
}

func TestGetRowShouldReturnImageColumn(t *testing.T) {
	image := mockTestGradientImage()

	lookupMap := make(map[color.Color]bool)
	expectedLength := mock_image_width

	row, err := GetImageRow(image, 0)
	assert.Nil(t, err)

	for _, color := range row {
		lookupMap[color] = true
	}

	acutalLength := len(lookupMap)
	assert.Equal(t, expectedLength, acutalLength)
}

func TestRotateImageShouldCorrectlyRotateImage(t *testing.T) {
	image := mockTestGradientImage()
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	imageRotated := RotateImage(image, 90)
	imageRotatedWidth := imageRotated.Bounds().Dx()
	imageRotatedHeight := imageRotated.Bounds().Dy()

	assert.Equal(t, imageWidth, imageRotatedWidth)
	assert.Equal(t, imageHeight, imageRotatedHeight)

	areNotEqual := false
	for y := 0; y < imageHeight; y += 1 {
		for x := 0; x < imageWidth; x += 1 {
			iR, iG, iB, _ := image.At(x, y).RGBA()
			irR, irG, irB, _ := imageRotated.At(x, y).RGBA()

			if iR != irR || iG != irG || iB != irB {
				areNotEqual = true
				break
			}
		}
	}
	assert.True(t, areNotEqual)

	imageRotated = RotateImage(imageRotated, 270)
	imageRotatedWidth = imageRotated.Bounds().Dx()
	imageRotatedHeight = imageRotated.Bounds().Dy()

	assert.Equal(t, imageWidth, imageRotatedWidth)
	assert.Equal(t, imageHeight, imageRotatedHeight)

	areEqual := true
	for y := 0; y < imageHeight; y += 1 {
		for x := 0; x < imageWidth; x += 1 {
			iR, iG, iB, _ := image.At(x, y).RGBA()
			irR, irG, irB, _ := imageRotated.At(x, y).RGBA()

			if iR != irR || iG != irG || iB != irB {
				areNotEqual = false
				break
			}
		}
	}
	assert.True(t, areEqual)
}

func TestTrimImageTransparentWorkspaceSohuldCorrectlyTrimRotatedImage(t *testing.T) {
	image := mockTestGradientImage()
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	imageRotated := RotateImage(image, 45)
	imageRotated = RotateImage(imageRotated, -45)
	imageRotatedWidth := imageRotated.Bounds().Dx()
	imageRotatedHeight := imageRotated.Bounds().Dy()

	assert.NotEqual(t, imageWidth, imageRotatedWidth)
	assert.NotEqual(t, imageHeight, imageRotatedHeight)

	imageTrimmed := TrimImageTransparentWorkspace(imageRotated, image)
	imageTrimmedWidth := imageTrimmed.Bounds().Dx()
	imageTrimmedHeight := imageTrimmed.Bounds().Dy()

	assert.Equal(t, imageWidth, imageTrimmedWidth)
	assert.Equal(t, imageHeight, imageTrimmedHeight)
}

func TestImageInvertShouldInvertImage(t *testing.T) {
	image := mockTestWhiteImage()
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	expectedColor := color.RGBA{0, 0, 0, 0xff}

	for y := 0; y < imageHeight; y += 1 {
		for x := 0; x < imageWidth; x += 1 {
			actualColor := image.At(x, y)

			assert.Equal(t, expectedColor, actualColor)
		}
	}
}

func TestImageShouldResize(t *testing.T) {
	image := mockTestWhiteImage()
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	expectedWidth := imageWidth / 2
	expectedHeight := imageHeight / 2

	scaledImage, err := ScaleImage(image, 0.5)
	assert.NotNil(t, scaledImage)
	assert.Nil(t, err)

	actualWidth := scaledImage.Bounds().Dx()
	actualHeight := scaledImage.Bounds().Dy()

	assert.Equal(t, expectedWidth, actualWidth)
	assert.Equal(t, expectedHeight, actualHeight)
}

const (
	mock_image_width  = 25
	mock_image_height = 25
)

// Create a test image which is a linear, left to right, black to white gradient of the size specifed by the mock_image prefixed constants
func mockTestGradientImage() draw.Image {
	gradient := make([]color.Color, mock_image_width)
	gradientStep := 255 / mock_image_width

	for gIndex := 0; gIndex < mock_image_width; gIndex += 1 {
		currentStep := gIndex * gradientStep
		if currentStep > 255 {
			currentStep = 255
		}

		gradient[gIndex] = color.RGBA{
			uint8(currentStep),
			uint8(currentStep),
			uint8(currentStep),
			0xff,
		}
	}

	p1 := image.Point{0, 0}
	p2 := image.Point{mock_image_width, mock_image_height}
	image := image.NewRGBA(image.Rectangle{p1, p2})

	for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
		for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
			image.Set(xIndex, yIndex, gradient[xIndex])
		}
	}

	return image
}

func mockTestWhiteImage() draw.Image {
	p1 := image.Point{0, 0}
	p2 := image.Point{mock_image_width, mock_image_height}
	image := image.NewRGBA(image.Rectangle{p1, p2})

	for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
		for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
			image.Set(xIndex, yIndex, color.Black)
		}
	}

	return image
}
