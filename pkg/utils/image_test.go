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

func TestSetColumnShouldSetColorsForTheImageColumn(t *testing.T) {
	image := mockTestWhiteImage()
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	black := color.RGBA{0x00, 0x00, 0x00, 0xff}
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}
	xIndexTarget := 2

	column := make([]color.Color, height)
	for yIndex := 0; yIndex < height; yIndex += 1 {
		column[yIndex] = color.Black
	}

	err := SetImageColumn(image, column, xIndexTarget)
	assert.Nil(t, err)

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			c := image.At(xIndex, yIndex)

			if xIndex == xIndexTarget {
				assert.Equal(t, black, c)
			} else {
				assert.Equal(t, white, c)
			}
		}
	}
}

func TestSetRowShouldSetColorsForTheImageRow(t *testing.T) {
	image := mockTestWhiteImage()
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	black := color.RGBA{0x00, 0x00, 0x00, 0xff}
	white := color.RGBA{0xff, 0xff, 0xff, 0xff}
	yIndexTarget := 2

	row := make([]color.Color, width)
	for xIndex := 0; xIndex < width; xIndex += 1 {
		row[xIndex] = color.Black
	}

	err := SetImageRow(image, row, yIndexTarget)
	assert.Nil(t, err)

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			c := image.At(xIndex, yIndex)

			if yIndex == yIndexTarget {
				assert.Equal(t, black, c)
			} else {
				assert.Equal(t, white, c)
			}
		}
	}
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

	invertedImage, err := InvertImage(image)

	assert.Nil(t, err)

	for y := 0; y < imageHeight; y += 1 {
		for x := 0; x < imageWidth; x += 1 {
			actualColor := invertedImage.At(x, y)

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

func TestShouldBlendImagesUsingLightenOnlyMode(t *testing.T) {
	cases := map[struct {
		a color.RGBA
		b color.RGBA
	}]color.RGBA{
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{0, 0, 0, 0xff}}:             {0, 0, 0, 0xff},
		{color.RGBA{255, 255, 255, 0xff}, color.RGBA{255, 255, 255, 0xff}}: {255, 255, 255, 0xff},
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{255, 255, 255, 0xff}}:       {255, 255, 255, 0xff},
		{color.RGBA{25, 50, 200, 0xff}, color.RGBA{200, 40, 20, 0xff}}:     {200, 50, 200, 0xff},
	}

	const height int = 2
	const width int = 2

	for c, expected := range cases {
		rect := image.Rect(0, 0, width, height)

		aImage := image.NewRGBA(rect)
		bImage := image.NewRGBA(rect)
		expectedImage := image.NewRGBA(rect)

		for xIndex := 0; xIndex < width; xIndex += 1 {
			for yIndex := 0; yIndex < height; yIndex += 1 {
				aImage.SetRGBA(xIndex, yIndex, c.a)
				bImage.SetRGBA(xIndex, yIndex, c.b)
				expectedImage.SetRGBA(xIndex, yIndex, expected)
			}
		}

		actualImage, err := BlendImages(aImage, bImage, LightenOnly)

		assert.Nil(t, err)
		assert.NotNil(t, actualImage)

		for xIndex := 0; xIndex < width; xIndex += 1 {
			for yIndex := 0; yIndex < height; yIndex += 1 {
				actualColor := actualImage.At(xIndex, yIndex)
				expectedColor := expectedImage.At(xIndex, yIndex)

				assert.Equal(t, expectedColor, actualColor)
			}
		}
	}
}

func TestShouldBlendImagesUsingDarkenOnlyMode(t *testing.T) {
	cases := map[struct {
		a color.RGBA
		b color.RGBA
	}]color.RGBA{
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{0, 0, 0, 0xff}}:             {0, 0, 0, 0xff},
		{color.RGBA{255, 255, 255, 0xff}, color.RGBA{255, 255, 255, 0xff}}: {255, 255, 255, 0xff},
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{255, 255, 255, 0xff}}:       {0, 0, 0, 0xff},
		{color.RGBA{25, 50, 200, 0xff}, color.RGBA{200, 40, 20, 0xff}}:     {25, 40, 20, 0xff},
	}

	const height int = 2
	const width int = 2

	for c, expected := range cases {
		rect := image.Rect(0, 0, width, height)

		aImage := image.NewRGBA(rect)
		bImage := image.NewRGBA(rect)
		expectedImage := image.NewRGBA(rect)

		for xIndex := 0; xIndex < width; xIndex += 1 {
			for yIndex := 0; yIndex < height; yIndex += 1 {
				aImage.SetRGBA(xIndex, yIndex, c.a)
				bImage.SetRGBA(xIndex, yIndex, c.b)
				expectedImage.SetRGBA(xIndex, yIndex, expected)
			}
		}

		actualImage, err := BlendImages(aImage, bImage, DarkenOnly)

		assert.Nil(t, err)
		assert.NotNil(t, actualImage)

		for xIndex := 0; xIndex < width; xIndex += 1 {
			for yIndex := 0; yIndex < height; yIndex += 1 {
				actualColor := actualImage.At(xIndex, yIndex)
				expectedColor := expectedImage.At(xIndex, yIndex)

				assert.Equal(t, expectedColor, actualColor)
			}
		}
	}
}

func TestShouldNotBlendImagesWithDifferentWidth(t *testing.T) {
	aImage := image.NewRGBA(image.Rect(0, 0, 2, 4))
	bImage := image.NewRGBA(image.Rect(0, 0, 4, 4))

	resultImage, err := BlendImages(aImage, bImage, LightenOnly)

	assert.NotNil(t, err)
	assert.Nil(t, resultImage)
}

func TestShouldNotBlendImagesWithDifferentHeight(t *testing.T) {
	aImage := image.NewRGBA(image.Rect(0, 0, 4, 2))
	bImage := image.NewRGBA(image.Rect(0, 0, 4, 4))

	resultImage, err := BlendImages(aImage, bImage, LightenOnly)

	assert.NotNil(t, err)
	assert.Nil(t, resultImage)
}

func TestAverageImageColumnShouldNotAverageForInvalidInput(t *testing.T) {
	cases := make([][][]color.Color, 0)

	cases = append(cases, [][]color.Color{})

	cases = append(cases, [][]color.Color{{}})

	cases = append(cases, [][]color.Color{
		{color.RGBA{0, 0, 0, 0xff}},
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{1, 1, 1, 0xff}},
	})

	for _, c := range cases {
		result, err := AverageImageColumns(c)

		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}

func TestAverageImageColumnShouldAverageSingleColumn(t *testing.T) {
	expectedColumn := []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{10, 10, 10, 0xff},
		color.RGBA{240, 240, 240, 0xff},
		color.RGBA{120, 120, 120, 0xff},
	}

	columns := make([][]color.Color, 0, 1)
	columns = append(columns, expectedColumn)
	actualColumn, err := AverageImageColumns(columns)

	assert.Nil(t, err)
	assert.NotNil(t, actualColumn)
	assert.Equal(t, expectedColumn, actualColumn)
}

func TestAverageImageColumnShouldAverageTwoColumns(t *testing.T) {
	columns := make([][]color.Color, 0, 2)

	columns = append(columns, []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{10, 10, 10, 0xff},
		color.RGBA{20, 20, 20, 0xff},
		color.RGBA{40, 40, 40, 0xff},
	})

	columns = append(columns, []color.Color{
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{200, 200, 200, 0xff},
		color.RGBA{200, 200, 200, 0xff},
	})

	expectedColumns := []color.Color{
		color.RGBA{52, 52, 52, 0xff},
		color.RGBA{52, 52, 52, 0xff},
		color.RGBA{115, 115, 115, 0xff},
		color.RGBA{115, 115, 115, 0xff},
	}

	actualColumn, err := AverageImageColumns(columns)

	assert.Nil(t, err)
	assert.NotNil(t, actualColumn)
	assert.Equal(t, expectedColumns, actualColumn)
}

func TestAverageImageColumnShouldAverageThreeColumns(t *testing.T) {
	columns := make([][]color.Color, 0, 3)

	columns = append(columns, []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{10, 10, 10, 0xff},
		color.RGBA{20, 20, 20, 0xff},
		color.RGBA{40, 40, 40, 0xff},
	})

	columns = append(columns, []color.Color{
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{200, 200, 200, 0xff},
		color.RGBA{200, 200, 200, 0xff},
	})

	columns = append(columns, []color.Color{
		color.RGBA{50, 50, 50, 0xff},
		color.RGBA{50, 50, 50, 0xff},
		color.RGBA{50, 50, 50, 0xff},
		color.RGBA{50, 50, 50, 0xff},
	})

	expectedColumns := []color.Color{
		color.RGBA{64, 64, 64, 0xff},
		color.RGBA{64, 64, 64, 0xff},
		color.RGBA{64, 64, 64, 0xff},
		color.RGBA{96, 96, 96, 0xff},
	}

	actualColumn, err := AverageImageColumns(columns)

	assert.Nil(t, err)
	assert.NotNil(t, actualColumn)
	assert.Equal(t, expectedColumns, actualColumn)
}

func TestAverageImageRowShouldNotAverageForInvalidInput(t *testing.T) {
	cases := make([][][]color.Color, 0)

	cases = append(cases, [][]color.Color{})

	cases = append(cases, [][]color.Color{{}})

	cases = append(cases, [][]color.Color{
		{color.RGBA{0, 0, 0, 0xff}},
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{1, 1, 1, 0xff}},
	})

	for _, c := range cases {
		result, err := AverageImageRow(c)

		assert.Nil(t, result)
		assert.NotNil(t, err)
	}
}

func TestAverageImageRowShouldAverageSingleRow(t *testing.T) {
	expectedRow := []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{10, 10, 10, 0xff},
		color.RGBA{240, 240, 240, 0xff},
		color.RGBA{120, 120, 120, 0xff},
	}

	rows := make([][]color.Color, 0, 1)
	rows = append(rows, expectedRow)
	actualRow, err := AverageImageRow(rows)

	assert.Nil(t, err)
	assert.NotNil(t, actualRow)
	assert.Equal(t, expectedRow, actualRow)
}

func TestAverageImageRowShouldAverageTwoRows(t *testing.T) {
	rows := make([][]color.Color, 0, 2)

	rows = append(rows, []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{10, 10, 10, 0xff},
		color.RGBA{20, 20, 20, 0xff},
		color.RGBA{40, 40, 40, 0xff},
	})

	rows = append(rows, []color.Color{
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{200, 200, 200, 0xff},
		color.RGBA{200, 200, 200, 0xff},
	})

	expectedRow := []color.Color{
		color.RGBA{52, 52, 52, 0xff},
		color.RGBA{52, 52, 52, 0xff},
		color.RGBA{115, 115, 115, 0xff},
		color.RGBA{115, 115, 115, 0xff},
	}

	actualRow, err := AverageImageRow(rows)

	assert.Nil(t, err)
	assert.NotNil(t, actualRow)
	assert.Equal(t, expectedRow, actualRow)
}

func TestAverageImageRowShouldAverageThreeRows(t *testing.T) {
	rows := make([][]color.Color, 0, 3)

	rows = append(rows, []color.Color{
		color.RGBA{0, 0, 0, 0xff},
		color.RGBA{10, 10, 10, 0xff},
		color.RGBA{20, 20, 20, 0xff},
		color.RGBA{40, 40, 40, 0xff},
	})

	rows = append(rows, []color.Color{
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{100, 100, 100, 0xff},
		color.RGBA{200, 200, 200, 0xff},
		color.RGBA{200, 200, 200, 0xff},
	})

	rows = append(rows, []color.Color{
		color.RGBA{50, 50, 50, 0xff},
		color.RGBA{50, 50, 50, 0xff},
		color.RGBA{50, 50, 50, 0xff},
		color.RGBA{50, 50, 50, 0xff},
	})

	expectedRow := []color.Color{
		color.RGBA{64, 64, 64, 0xff},
		color.RGBA{64, 64, 64, 0xff},
		color.RGBA{64, 64, 64, 0xff},
		color.RGBA{96, 96, 96, 0xff},
	}

	acutalRow, err := AverageImageRow(rows)

	assert.Nil(t, err)
	assert.NotNil(t, acutalRow)
	assert.Equal(t, expectedRow, acutalRow)
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
			image.Set(xIndex, yIndex, color.White)
		}
	}

	return image
}
