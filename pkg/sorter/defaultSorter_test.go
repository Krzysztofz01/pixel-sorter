package sorter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSorterShouldSortForSortByBrightnessSplitByBrightnessZeroAngle(t *testing.T) {
	sorterOptions := SorterOptions{
		SortByBrightness,
		SplitByBrightness,
		0.0,
		1.0,
		0,
	}

	// NOTE: Preparing the expected image
	p1 := image.Point{0, 0}
	p2 := image.Point{mock_image_width, mock_image_height}
	cBlack := color.RGBA{0, 0, 0, 0xff}
	cWhite := color.RGBA{255, 255, 255, 0xff}
	expectedImage := image.NewRGBA(image.Rectangle{p1, p2})
	for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
		for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
			if xIndex < 2 {
				expectedImage.Set(xIndex, yIndex, cBlack)
			} else {
				expectedImage.Set(xIndex, yIndex, cWhite)
			}
		}
	}

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), &sorterOptions)
	assert.Nil(t, err)

	actualImage, err := sorter.Sort()
	assert.Nil(t, err)

	for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
		for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
			eR, eG, eB, _ := expectedImage.At(xIndex, yIndex).RGBA()
			aR, aG, aB, _ := actualImage.At(xIndex, yIndex).RGBA()

			assert.Equal(t, eR, aR)
			assert.Equal(t, eG, aG)
			assert.Equal(t, eB, aB)
		}
	}
}

const (
	mock_image_width  = 5
	mock_image_height = 5
)

// Create a test image which consists of black and white 1px wide columns of the size specified by the mock_image prefixed constants
func mockTestBlackAndWhiteStripesImage() draw.Image {
	p1 := image.Point{0, 0}
	p2 := image.Point{mock_image_width, mock_image_height}
	image := image.NewRGBA(image.Rectangle{p1, p2})

	cBlack := color.RGBA{0, 0, 0, 0xff}
	cWhite := color.RGBA{255, 255, 255, 0xff}

	for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
		for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
			if xIndex%2 != 0 {
				image.Set(xIndex, yIndex, cBlack)
			} else {
				image.Set(xIndex, yIndex, cWhite)
			}
		}
	}

	return image
}
