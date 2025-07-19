package img

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/Krzysztofz01/imaging"
	"github.com/Krzysztofz01/pimit"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

// TODO: Fine tune const values in order to make the edge detection more versatile
// TODO: Fix the NonMaxSuppresion? to prevent the creation of additional noise on the image
// TODO: Add edge detection unit tests

const (
	blurSigmaParam            float64 = 1.3
	doubleThresholdZero       uint8   = 0
	doubleThresholdWeak       uint8   = 75
	doubleThresholdStrong     uint8   = 255
	doubleThresholdLowerRatio float64 = 0.05
	doubleThresholdUpperRatio float64 = 0.09
	lowerHysteresisThreshold  float64 = 0.25
	upperHysteresisThreshold  float64 = 0.75
)

var (
	sobelMatrixVertical [9]float64 = [9]float64{
		-1, -2, -1,
		0, 0, 0,
		1, 2, 1,
	}

	sobelMatrixHorizontal [9]float64 = [9]float64{
		-1, 0, 1,
		-2, 0, 2,
		-1, 0, 1,
	}
)

// Helper structure used to represent gradient point magnitude and angle direction
type gradientPoint struct {
	magnitude float64
	direction float64
}

// Generate a edge detection image based on the given input image using the Canny edge detection algorithm
func PerformEdgeDetection(i *image.NRGBA, performNonMaxSupression, invertColors bool) (*image.NRGBA, error) {
	imgGrayscale := utils.GrayscaleNrgba(i)

	imgSmoothed := imaging.Blur(imgGrayscale, blurSigmaParam)

	imgVerticalConv := imaging.Convolve3x3(imgSmoothed, sobelMatrixVertical, nil)

	imgHorizontalConv := imaging.Convolve3x3(imgSmoothed, sobelMatrixHorizontal, nil)

	gradientPoints, err := calculateGradientPoints(imgVerticalConv, imgHorizontalConv)
	if err != nil {
		return nil, fmt.Errorf("edge-detection: failed to calculate the gradient points: %w", err)
	}

	var imgGradient *image.NRGBA = nil
	if performNonMaxSupression {
		imgGradient, err = performNonMaxSuppresion(imgSmoothed, gradientPoints)
		if err != nil {
			return nil, fmt.Errorf("edge-detection: failed to perform no max suppresion on the image: %w", err)
		}
	} else {
		imgGradient, err = createGradientMapImage(imgSmoothed, gradientPoints)
		if err != nil {
			return nil, fmt.Errorf("edge-detection: failed to generate the gradient points image: %w", err)
		}
	}

	imgHysteresis := performHysteresis(imgGradient)
	if invertColors {
		return utils.InvertImageNrgba(imgHysteresis), nil
	} else {
		return imgHysteresis, nil
	}
}

// Helper function used to create a gradient point matrix by calculating the sobel vertical and horizontal derivatives
func calculateGradientPoints(vertical *image.NRGBA, horizontal *image.NRGBA) ([][]gradientPoint, error) {
	width := vertical.Bounds().Dx()
	height := vertical.Bounds().Dy()
	magnitudeMaxValue := math.Inf(-1)
	magnitudeMaxLock := sync.Mutex{}

	gp := make([][]gradientPoint, width)
	for x := 0; x < width; x += 1 {
		gp[x] = make([]gradientPoint, height)
	}

	pimit.ParallelMatrixReadWrite(gp, func(xIndex, yIndex int, _ gradientPoint) gradientPoint {
		vColor := float64(utils.NrgbaToGrayscaleComponent(vertical.NRGBAAt(xIndex, yIndex)))
		hColor := float64(utils.NrgbaToGrayscaleComponent(horizontal.NRGBAAt(xIndex, yIndex)))

		magnitude := math.Hypot(vColor, hColor)
		direction := math.Atan2(vColor, hColor)

		// TODO: The changes of magnitudeMax caused data-race in v0.2.0-beta release. This can be fixed in the
		// future using a atomic float64. For now, we are locking using a mutex before magnitudeMax access.
		// (We can even remove it!)
		magnitudeMaxLock.Lock()
		if magnitude > magnitudeMaxValue {
			magnitudeMaxValue = magnitude
		}
		magnitudeMaxLock.Unlock()

		return gradientPoint{
			magnitude: magnitude,
			direction: direction,
		}
	})

	return gp, nil
}

// Helper function used to generate a NRGBA image from the gradient points values
func createGradientMapImage(i *image.NRGBA, g [][]gradientPoint) (*image.NRGBA, error) {
	if ok := validateGradientPointsWithImage(i, g); !ok {
		return nil, errors.New("edge-detection: the provided image and gradient points are not correspoding")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	outputImage := newWhiteNRGBA(image.Rect(0, 0, width, height))
	pimit.ParallelReadWrite(outputImage, func(xIndex, yIndex int, _ color.Color) color.Color {
		magnitude := g[xIndex][yIndex].magnitude
		colorValue := uint8(utils.ClampFloat64(0, magnitude, 255))

		return color.Gray{Y: colorValue}
	})

	return outputImage, nil
}

// Helper function used to perform a Non Max Supression operation on the image
func performNonMaxSuppresion(i *image.NRGBA, g [][]gradientPoint) (*image.NRGBA, error) {
	if ok := validateGradientPointsWithImage(i, g); !ok {
		return nil, errors.New("edge-detection: the provided image and gradient points are not correspoding")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	outputImage := newWhiteNRGBA(image.Rect(0, 0, width, height))
	pimit.ParallelReadWrite(outputImage, func(xIndex, yIndex int, _ color.Color) color.Color {
		magnitude := g[xIndex][yIndex].magnitude
		colorValue := uint8(utils.ClampFloat64(0, magnitude, 255))

		if xIndex == 0 || xIndex == width-1 || yIndex == 0 || yIndex == height-1 {
			return color.Gray{Y: colorValue}
		}

		magnitudeIsMax := false
		direction := g[xIndex][yIndex].direction * 180.0 / math.Pi
		if direction < 0 {
			direction += 180.0
		}

		if valueBetween(direction, 0, 22.5) || valueBetween(direction, 157.5, 180) {
			if magnitude > g[xIndex][yIndex+1].magnitude && magnitude > g[xIndex][yIndex-1].magnitude {
				magnitudeIsMax = true
			}
		}

		if valueBetween(direction, 22.5, 67.5) {
			if magnitude > g[xIndex+1][yIndex-1].magnitude && magnitude > g[xIndex-1][yIndex+1].magnitude {
				magnitudeIsMax = true
			}
		}

		if valueBetween(direction, 67.5, 112.5) {
			if magnitude > g[xIndex+1][yIndex].magnitude && magnitude > g[xIndex-1][yIndex].magnitude {
				magnitudeIsMax = true
			}
		}

		if valueBetween(direction, 112.5, 157.5) {
			if magnitude > g[xIndex-1][yIndex-1].magnitude && magnitude > g[xIndex+1][yIndex+1].magnitude {
				magnitudeIsMax = true
			}
		}

		if magnitudeIsMax {
			return color.White

		}

		return color.Gray{Y: colorValue}
	})

	return outputImage, nil
}

// Helper function used to perform a Double Threshold and Hysteresis operations on the image
func performHysteresis(i *image.NRGBA) *image.NRGBA {
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	maxColorValue := 0
	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			colorValue := utils.NrgbaToGrayscaleComponent(i.NRGBAAt(xIndex, yIndex))
			if colorValue > maxColorValue {
				maxColorValue = colorValue
			}
		}
	}

	dtUpperThreshold := float64(maxColorValue) * doubleThresholdUpperRatio
	dtLowerThreshold := dtUpperThreshold * doubleThresholdLowerRatio

	colorZero := color.Gray{Y: doubleThresholdZero}
	colorWeak := color.Gray{Y: doubleThresholdWeak}
	colorStrong := color.Gray{Y: doubleThresholdStrong}

	outputImage := newWhiteNRGBA(image.Rect(0, 0, width, height))
	pimit.ParallelReadWrite(outputImage, func(xIndex, yIndex int, _ color.Color) color.Color {
		colorValue := float64(utils.NrgbaToGrayscaleComponent(i.NRGBAAt(xIndex, yIndex)))

		if colorValue >= dtUpperThreshold {
			return colorStrong
		}

		if colorValue < dtLowerThreshold {
			return colorZero
		}

		return colorWeak
	})

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			if xIndex == 0 || xIndex == width-1 || yIndex == 0 || yIndex == height-1 {
				outputImage.Set(xIndex, yIndex, colorZero)
				continue
			}

			colorValue := uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex, yIndex)))
			if colorValue == doubleThresholdWeak {
				neighbors := [8]uint8{
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex-1, yIndex-1))),
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex-1, yIndex))),
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex-1, yIndex+1))),
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex, yIndex-1))),
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex, yIndex+1))),
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex+1, yIndex-1))),
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex+1, yIndex))),
					uint8(utils.NrgbaToGrayscaleComponent(outputImage.NRGBAAt(xIndex+1, yIndex+1))),
				}

				isColorStrong := false
				for _, neighborValue := range neighbors {
					if neighborValue == doubleThresholdStrong {
						isColorStrong = true
						break
					}
				}

				if isColorStrong {
					outputImage.Set(xIndex, yIndex, colorStrong)
				} else {
					outputImage.Set(xIndex, yIndex, colorZero)
				}
			}
		}
	}

	return outputImage
}

// Helper function used to check if the given values is between two edge values
func valueBetween(value, min, max float64) bool {
	if min > max {
		panic("edge-detection: invalid value between operation boundaries")
	}

	return value > min && value < max
}

// Helper function used to createa new image.NRGBA image initialize with 0xFFFFFFFF color
func newWhiteNRGBA(r image.Rectangle) *image.NRGBA {
	img := image.NewNRGBA(r)
	pimit.ParallelNrgbaReadWrite(img, func(_, _ int, _, _, _, _ uint8) (uint8, uint8, uint8, uint8) {
		return 255, 255, 255, 255
	})

	return img
}

// Helper function used to perform a validation of corresponding image and gradient points collection
func validateGradientPointsWithImage(i *image.NRGBA, g [][]gradientPoint) bool {
	if i == nil || g == nil {
		return false
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	if len(g) != width {
		return false
	}

	for _, gr := range g {
		if len(gr) != height {
			return false
		}
	}

	return true
}
