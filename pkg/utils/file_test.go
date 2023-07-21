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
	test_file_name_jpg string = "test-file-utility-image.jpg"
	test_file_name_png string = "test-file-utility-image.png"
)

func TestImageFileShouldBeStoredAsJpg(t *testing.T) {
	clearEnvironmentFromTestFiles()

	expectedImage := mockTestBlackImage()

	err := StoreImageToFile(test_file_name_jpg, "jpg", expectedImage)
	assert.Nil(t, err)

	file, err := os.Open(test_file_name_jpg)
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

	err := StoreImageToFile(test_file_name_png, "png", expectedImage)
	assert.Nil(t, err)

	file, err := os.Open(test_file_name_png)
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

	err := StoreImageToFile(test_file_name_png, "png", expectedImage)
	assert.Nil(t, err)

	actualImage, err := GetImageFromFile(test_file_name_png)
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

func TestEscapedPathQuotesShouldCorrectlyRemoveSurroundingQuotes(t *testing.T) {
	cases := map[string]struct {
		path string
		ok   bool
	}{
		"hello/world.jpg":                       {"hello/world.jpg", true},
		"hello/world.pNg":                       {"hello/world.pNg", true},
		"hello/world.gif":                       {"hello/world.gif", true},
		"'hello/world.jpg'":                     {"hello/world.jpg", true},
		"\"hello/world.png\"":                   {"hello/world.png", true},
		"\"'\"'hello/world.jpg'\"'\"":           {"hello/world.jpg", true},
		"'''''''''''hello/world.jpg'''''''''''": {"", false},
	}

	for path, expected := range cases {
		actualPath, err := EscapePathQuotes(path)

		assert.Equal(t, expected.path, actualPath)

		if expected.ok {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

// Helper function that is removing the image files created during the tests
func clearEnvironmentFromTestFiles() {
	if info, err := os.Stat(test_file_name_jpg); err == nil && !info.IsDir() {
		if err := os.Remove(test_file_name_jpg); err != nil {
			fmt.Println(err)
			panic("utils-test: can not remove jpg test file")
		}
	}

	if info, err := os.Stat(test_file_name_png); err == nil && !info.IsDir() {
		if err := os.Remove(test_file_name_png); err != nil {
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
