package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	test_file_name string = "test-file-utility-image"
)

func TestImageFileShouldBeStoredAsJpg(t *testing.T) {
	clearEnvironmentFromTestFiles()

	expectedImage := mockTestBlackImage()

	err := StoreImageToFile(test_file_name, "jpg", expectedImage)
	assert.Nil(t, err)

	filePath := fmt.Sprintf("%s-sorted.jpg", test_file_name)
	file, err := os.Open(filePath)
	assert.Nil(t, err)

	actualImage, err := jpeg.Decode(file)
	assert.Nil(t, err)

	err = file.Close()
	assert.Nil(t, err)

	assert.Equal(t, expectedImage.Bounds().Dx(), actualImage.Bounds().Dx())
	assert.Equal(t, expectedImage.Bounds().Dy(), actualImage.Bounds().Dy())

	for xIndex := 0; xIndex < expectedImage.Bounds().Dx(); xIndex += 1 {
		for yIndex := 0; yIndex < expectedImage.Bounds().Dy(); yIndex += 1 {
			eR, eG, eB, _ := expectedImage.At(xIndex, yIndex).RGBA()
			aR, aG, aB, _ := actualImage.At(xIndex, yIndex).RGBA()

			assert.Equal(t, eR, aR)
			assert.Equal(t, eG, aG)
			assert.Equal(t, eB, aB)
		}
	}

	clearEnvironmentFromTestFiles()
}

func TestImageFileShouldBeStoredAsPng(t *testing.T) {
	clearEnvironmentFromTestFiles()

	expectedImage := mockTestBlackImage()

	err := StoreImageToFile(test_file_name, "png", expectedImage)
	assert.Nil(t, err)

	filePath := fmt.Sprintf("%s-sorted.png", test_file_name)
	file, err := os.Open(filePath)
	assert.Nil(t, err)

	actualImage, err := png.Decode(file)
	assert.Nil(t, err)

	err = file.Close()
	assert.Nil(t, err)

	assert.Equal(t, expectedImage.Bounds().Dx(), actualImage.Bounds().Dx())
	assert.Equal(t, expectedImage.Bounds().Dy(), actualImage.Bounds().Dy())

	for xIndex := 0; xIndex < expectedImage.Bounds().Dx(); xIndex += 1 {
		for yIndex := 0; yIndex < expectedImage.Bounds().Dy(); yIndex += 1 {
			eR, eG, eB, _ := expectedImage.At(xIndex, yIndex).RGBA()
			aR, aG, aB, _ := actualImage.At(xIndex, yIndex).RGBA()

			assert.Equal(t, eR, aR)
			assert.Equal(t, eG, aG)
			assert.Equal(t, eB, aB)
		}
	}

	clearEnvironmentFromTestFiles()
}

func TestImageShouldBeRetrievied(t *testing.T) {
	clearEnvironmentFromTestFiles()

	expectedImage := mockTestBlackImage()

	err := StoreImageToFile(test_file_name, "png", expectedImage)
	assert.Nil(t, err)

	actualImageFilePath := fmt.Sprintf("%s-sorted.png", test_file_name)
	actualImage, err := GetImageFromFile(actualImageFilePath)
	assert.Nil(t, err)

	assert.Equal(t, expectedImage.Bounds().Dx(), actualImage.Bounds().Dx())
	assert.Equal(t, expectedImage.Bounds().Dy(), actualImage.Bounds().Dy())

	for xIndex := 0; xIndex < expectedImage.Bounds().Dx(); xIndex += 1 {
		for yIndex := 0; yIndex < expectedImage.Bounds().Dy(); yIndex += 1 {
			eR, eG, eB, _ := expectedImage.At(xIndex, yIndex).RGBA()
			aR, aG, aB, _ := actualImage.At(xIndex, yIndex).RGBA()

			assert.Equal(t, eR, aR)
			assert.Equal(t, eG, aG)
			assert.Equal(t, eB, aB)
		}
	}

	clearEnvironmentFromTestFiles()
}

// Helper function that is removing the image files created during the tests
func clearEnvironmentFromTestFiles() {
	jpgFileName := fmt.Sprintf("%s-sorted.jpg", test_file_name)
	if info, err := os.Stat(jpgFileName); err == nil && !info.IsDir() {
		if err := os.Remove(jpgFileName); err != nil {
			fmt.Println(err)
			panic("utils-test: can not remove jpg test file")
		}
	}

	pngFileName := fmt.Sprintf("%s-sorted.png", test_file_name)
	if info, err := os.Stat(pngFileName); err == nil && !info.IsDir() {
		if err := os.Remove(pngFileName); err != nil {
			panic("utils-test: can not remove jpg test file")
		}
	}
}

// Create a test image which is a linear, left to right, black to white gradient of the size specifed by the mock_image prefixed constants
func mockTestBlackImage() image.Image {
	const width int = 10
	const height int = 10

	image := image.NewRGBA(image.Rect(0, 0, width, height))
	color := color.RGBA{0, 0, 0, 0xff}

	for yIndex := 0; yIndex < height; yIndex += 1 {
		for xIndex := 0; xIndex < width; xIndex += 1 {
			image.Set(xIndex, yIndex, color)
		}
	}

	return image
}
