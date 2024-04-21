package utils

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestImageToNrgbaShouldConvertImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := map[struct {
		color color.Color
		image draw.Image
	}]color.NRGBA{
		{color.RGBA{255, 255, 255, 255}, image.NewRGBA(image.Rect(0, 0, 10, 10))}:   {255, 255, 255, 255},
		{color.NRGBA{255, 255, 255, 255}, image.NewNRGBA(image.Rect(0, 0, 10, 10))}: {255, 255, 255, 255},
		{color.Gray{255}, image.NewGray(image.Rect(0, 0, 10, 10))}:                  {255, 255, 255, 255},
	}

	for c, expectedColor := range cases {
		for y := 0; y < c.image.Bounds().Dy(); y += 1 {
			for x := 0; x < c.image.Bounds().Dx(); x += 1 {
				c.image.Set(x, y, c.color)
			}
		}

		convertedImage := ImageToNrgbaImage(c.image)
		assert.NotNil(t, convertedImage)

		for y := 0; y < convertedImage.Rect.Dy(); y += 1 {
			for x := 0; x < convertedImage.Rect.Dx(); x += 1 {
				assert.Equal(t, expectedColor, convertedImage.At(x, y))
			}
		}
	}
}

func TestImageToRgbaShouldConvertImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := map[struct {
		color color.Color
		image draw.Image
	}]color.RGBA{
		{color.RGBA{255, 255, 255, 255}, image.NewRGBA(image.Rect(0, 0, 10, 10))}:   {255, 255, 255, 255},
		{color.NRGBA{255, 255, 255, 255}, image.NewNRGBA(image.Rect(0, 0, 10, 10))}: {255, 255, 255, 255},
		{color.Gray{255}, image.NewGray(image.Rect(0, 0, 10, 10))}:                  {255, 255, 255, 255},
	}

	for c, expectedColor := range cases {
		for y := 0; y < c.image.Bounds().Dy(); y += 1 {
			for x := 0; x < c.image.Bounds().Dx(); x += 1 {
				c.image.Set(x, y, c.color)
			}
		}

		convertedImage := ImageToRgbaImage(c.image)
		assert.NotNil(t, convertedImage)

		for y := 0; y < convertedImage.Rect.Dy(); y += 1 {
			for x := 0; x < convertedImage.Rect.Dx(); x += 1 {
				assert.Equal(t, expectedColor, convertedImage.At(x, y))
			}
		}
	}
}

func TestGetImageCopyNrgbaShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		GetImageCopyNrgba(nil)
	})
}

func TestGetImageCopyNrgbaShouldCreateACopy(t *testing.T) {
	defer goleak.VerifyNone(t)

	originalImage := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < originalImage.Rect.Dy(); y += 1 {
		for x := 0; x < originalImage.Rect.Dx(); x += 1 {
			originalImage.Set(x, y, color.White)
		}
	}

	copyImage := GetImageCopyNrgba(originalImage)
	assert.NotNil(t, copyImage)

	for y := 0; y < originalImage.Rect.Dy(); y += 1 {
		for x := 0; x < originalImage.Rect.Dx(); x += 1 {
			assert.Equal(t, copyImage.At(x, y), originalImage.At(x, y))

			originalImage.Set(x, y, color.Black)

			assert.NotEqual(t, copyImage.At(x, y), originalImage.At(x, y))
		}
	}
}

func TestGetImageCopyRgbaShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		GetImageCopyRgba(nil)
	})
}

func TestGetImageCopyRgbaShouldCreateACopy(t *testing.T) {
	defer goleak.VerifyNone(t)

	originalImage := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < originalImage.Rect.Dy(); y += 1 {
		for x := 0; x < originalImage.Rect.Dx(); x += 1 {
			originalImage.Set(x, y, color.White)
		}
	}

	copyImage := GetImageCopyRgba(originalImage)
	assert.NotNil(t, copyImage)

	for y := 0; y < originalImage.Rect.Dy(); y += 1 {
		for x := 0; x < originalImage.Rect.Dx(); x += 1 {
			assert.Equal(t, copyImage.At(x, y), originalImage.At(x, y))

			originalImage.Set(x, y, color.Black)

			assert.NotEqual(t, copyImage.At(x, y), originalImage.At(x, y))
		}
	}
}

func TestInvertImageNrgbaShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		InvertImageNrgba(nil)
	})
}

func TestInvertImageNrgbaShouldInvertImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	originalImage := image.NewNRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < originalImage.Rect.Dy(); y += 1 {
		for x := 0; x < originalImage.Rect.Dx(); x += 1 {
			originalImage.Set(x, y, color.White)
		}
	}

	expectedColor := color.NRGBA{0, 0, 0, 0xff}

	invertedImage := InvertImageNrgba(originalImage)

	assert.NotNil(t, invertedImage)

	for y := 0; y < originalImage.Rect.Dy(); y += 1 {
		for x := 0; x < originalImage.Rect.Dx(); x += 1 {
			assert.Equal(t, expectedColor, invertedImage.At(x, y))
		}
	}
}

func TestRotateImageNrgbaShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		RotateImageNrgba(nil, 90)
	})
}

func TestRotateImageNrgbaShouldCorrectlyRotateImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockTestGradientImageNrgba()
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	imageRotated := RotateImageNrgba(image, 90)
	imageRotatedWidth := imageRotated.Bounds().Dx()
	imageRotatedHeight := imageRotated.Bounds().Dy()

	assert.Equal(t, imageWidth, imageRotatedWidth)
	assert.Equal(t, imageHeight, imageRotatedHeight)

	// NOTE: Checking just the first half of the first row
	for x := 0; x < imageWidth/2; x += 1 {
		assert.NotEqual(t, image.NRGBAAt(x, 0), imageRotated.NRGBAAt(x, 0))
	}

	imageRotated = RotateImageNrgba(imageRotated, 270)
	imageRotatedWidth = imageRotated.Bounds().Dx()
	imageRotatedHeight = imageRotated.Bounds().Dy()

	assert.Equal(t, imageWidth, imageRotatedWidth)
	assert.Equal(t, imageHeight, imageRotatedHeight)

	for y := 0; y < imageHeight; y += 1 {
		for x := 0; x < imageWidth; x += 1 {
			assert.Equal(t, image.NRGBAAt(x, y), imageRotated.NRGBAAt(x, y))
		}
	}
}

func TestRotateImageNrgbaShouldCorrectlyRotateImageWithFullRotation(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := []int{0, 360}

	for _, angle := range cases {
		image := mockTestGradientImageNrgba()
		imageWidth := image.Bounds().Dx()
		imageHeight := image.Bounds().Dy()

		imageRotated := RotateImageNrgba(image, angle)
		imageRotatedWidth := imageRotated.Bounds().Dx()
		imageRotatedHeight := imageRotated.Bounds().Dy()

		assert.Equal(t, imageWidth, imageRotatedWidth)
		assert.Equal(t, imageHeight, imageRotatedHeight)

		for y := 0; y < imageHeight; y += 1 {
			for x := 0; x < imageWidth; x += 1 {
				assert.Equal(t, image.NRGBAAt(x, y), imageRotated.NRGBAAt(x, y))
			}
		}
	}
}

func TestRotateImageWithRevertNrgbaShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		RotateImageWithRevertNrgba(nil, 90)
	})
}

func TestRotateImageWithRevertNrgbaShouldRotateAndRevert(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockTestGradientImageNrgba()
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	imageRotated, revert := RotateImageWithRevertNrgba(image, 45)
	imageRotatedWidth := imageRotated.Bounds().Dx()
	imageRotatedHeight := imageRotated.Bounds().Dy()

	assert.NotEqual(t, imageWidth, imageRotatedWidth)
	assert.NotEqual(t, imageHeight, imageRotatedHeight)

	imageRevert := revert(imageRotated)
	imageRevertWidth := imageRevert.Bounds().Dx()
	imageRevertHeight := imageRevert.Bounds().Dy()

	assert.Equal(t, imageWidth, imageRevertWidth)
	assert.Equal(t, imageHeight, imageRevertHeight)
}

func TestRotateImageWithRevertNrgbaShouldRotateAndRevertWithFullRotation(t *testing.T) {
	defer goleak.VerifyNone(t)

	image := mockTestGradientImageNrgba()
	imageWidth := image.Bounds().Dx()
	imageHeight := image.Bounds().Dy()

	imageRotated, revert := RotateImageWithRevertNrgba(image, 0)
	imageRotatedWidth := imageRotated.Bounds().Dx()
	imageRotatedHeight := imageRotated.Bounds().Dy()

	assert.Equal(t, imageWidth, imageRotatedWidth)
	assert.Equal(t, imageHeight, imageRotatedHeight)

	imageRevert := revert(imageRotated)
	imageRevertWidth := imageRevert.Bounds().Dx()
	imageRevertHeight := imageRevert.Bounds().Dy()

	assert.Equal(t, imageWidth, imageRevertWidth)
	assert.Equal(t, imageHeight, imageRevertHeight)
}

func TestRotateWithRevertNrgbaShouldPanicOnInvalidRevertArguments(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := mockTestGradientImageNrgba()
	imageWidth := img.Bounds().Dx()
	imageHeight := img.Bounds().Dy()

	_, revert := RotateImageWithRevertNrgba(img, 0)

	assert.Panics(t, func() {
		revert(nil)
	})

	assert.Panics(t, func() {
		invalidImage := image.NewNRGBA(image.Rect(0, 0, imageWidth*2, imageHeight*3))
		revert(invalidImage)
	})
}

func TestScaleImageNrgbaShouldPanicOnNilImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	assert.Panics(t, func() {
		ScaleImageNrgba(nil, 0.5)
	})
}

func TestScaleImageNrgbaShouldErrorOnInvalidPercentage(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := image.NewNRGBA(image.Rect(0, 0, 10, 10))

	_, err := ScaleImageNrgba(img, 0.0)
	assert.NotNil(t, err)

	_, err = ScaleImageNrgba(img, 1.1)
	assert.NotNil(t, err)

	_, err = ScaleImageNrgba(img, -0.1)
	assert.NotNil(t, err)
}

func TestScaleImageNrgbaShouldResizeImage(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := []struct {
		width  int
		height int
		factor float64
	}{
		{10, 20, 0.5},
		{10, 20, 1},
	}

	for _, c := range cases {
		originalImage := image.NewNRGBA(image.Rect(0, 0, c.width, c.height))
		for y := 0; y < c.height; y += 1 {
			for x := 0; x < c.width; x += 1 {
				originalImage.Set(x, y, color.White)
			}
		}

		expectedWidth := int(float64(c.width) * c.factor)
		expectedHeight := int(float64(c.height) * c.factor)

		scaledImage, err := ScaleImageNrgba(originalImage, c.factor)
		assert.NotNil(t, scaledImage)
		assert.Nil(t, err)

		actualWidth := scaledImage.Bounds().Dx()
		actualHeight := scaledImage.Bounds().Dy()

		assert.Equal(t, expectedWidth, actualWidth)
		assert.Equal(t, expectedHeight, actualHeight)
	}
}

func TestBlendImagesNrgbaShouldBlendImagesUsingDarkenOnlyMode(t *testing.T) {
	defer goleak.VerifyNone(t)

	cases := map[struct {
		a color.NRGBA
		b color.NRGBA
	}]color.NRGBA{
		{color.NRGBA{0, 0, 0, 0xff}, color.NRGBA{0, 0, 0, 0xff}}:             {0, 0, 0, 0xff},
		{color.NRGBA{255, 255, 255, 0xff}, color.NRGBA{255, 255, 255, 0xff}}: {255, 255, 255, 0xff},
		{color.NRGBA{0, 0, 0, 0xff}, color.NRGBA{255, 255, 255, 0xff}}:       {0, 0, 0, 0xff},
		{color.NRGBA{25, 50, 200, 0xff}, color.NRGBA{200, 40, 20, 0xff}}:     {25, 40, 20, 0xff},
	}

	const height int = 2
	const width int = 2

	for c, expected := range cases {
		rect := image.Rect(0, 0, width, height)

		aImage := image.NewNRGBA(rect)
		bImage := image.NewNRGBA(rect)
		expectedImage := image.NewNRGBA(rect)

		for xIndex := 0; xIndex < width; xIndex += 1 {
			for yIndex := 0; yIndex < height; yIndex += 1 {
				aImage.SetNRGBA(xIndex, yIndex, c.a)
				bImage.SetNRGBA(xIndex, yIndex, c.b)
				expectedImage.SetNRGBA(xIndex, yIndex, expected)
			}
		}

		actualImage, err := BlendImagesNrgba(aImage, bImage, DarkenOnly)

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

func TestBlendImagesNrgbaShouldPanicOnNilImages(t *testing.T) {
	defer goleak.VerifyNone(t)

	img := image.NewNRGBA(image.Rect(0, 0, 10, 10))

	assert.Panics(t, func() {
		BlendImagesNrgba(nil, img, LightenOnly)
	})

	assert.Panics(t, func() {
		BlendImagesNrgba(img, nil, LightenOnly)
	})

	assert.Panics(t, func() {
		BlendImagesNrgba(nil, nil, LightenOnly)
	})
}

func TestBlendImagesNrgbaShouldNotBlendImagesWithDifferentWidth(t *testing.T) {
	defer goleak.VerifyNone(t)

	aImage := image.NewNRGBA(image.Rect(0, 0, 2, 4))
	bImage := image.NewNRGBA(image.Rect(0, 0, 4, 4))

	resultImage, err := BlendImagesNrgba(aImage, bImage, LightenOnly)

	assert.NotNil(t, err)
	assert.Nil(t, resultImage)
}

func TestBlendImagesNrgbaShouldNotBlendImagesWithDifferentHeight(t *testing.T) {
	defer goleak.VerifyNone(t)

	aImage := image.NewNRGBA(image.Rect(0, 0, 4, 2))
	bImage := image.NewNRGBA(image.Rect(0, 0, 4, 4))

	resultImage, err := BlendImagesNrgba(aImage, bImage, LightenOnly)

	assert.NotNil(t, err)
	assert.Nil(t, resultImage)
}

func TestGrayscaleShouldTurnImageIntoGrayscale(t *testing.T) {
	p1 := image.Point{0, 0}
	p2 := image.Point{mock_image_width, mock_image_height}
	image := image.NewNRGBA(image.Rectangle{p1, p2})

	for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
		for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
			image.SetNRGBA(xIndex, yIndex, color.NRGBA{0x2f, 0x0f, 0xe2, 0xff})
		}
	}

	grayscale := GrayscaleNrgba(image)
	for yIndex := 0; yIndex < grayscale.Rect.Dy(); yIndex += 1 {
		for xIndex := 0; xIndex < grayscale.Rect.Dx(); xIndex += 1 {
			r, g, b, _ := grayscale.At(xIndex, yIndex).RGBA()

			assert.True(t, r == g)
			assert.True(t, b == g)
			assert.True(t, r == b)
		}
	}
}

const (
	mock_image_width  = 25
	mock_image_height = 25
)

func mockTestGradientImageNrgba() *image.NRGBA {
	gradient := make([]color.NRGBA, mock_image_width)
	gradientStep := 255 / mock_image_width

	for gIndex := 0; gIndex < mock_image_width; gIndex += 1 {
		currentStep := gIndex * gradientStep
		if currentStep > 255 {
			currentStep = 255
		}

		gradient[gIndex] = color.NRGBA{
			uint8(currentStep),
			uint8(currentStep),
			uint8(currentStep),
			0xff,
		}
	}

	p1 := image.Point{0, 0}
	p2 := image.Point{mock_image_width, mock_image_height}
	image := image.NewNRGBA(image.Rectangle{p1, p2})

	for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
		for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
			image.SetNRGBA(xIndex, yIndex, gradient[xIndex])
		}
	}

	return image
}
