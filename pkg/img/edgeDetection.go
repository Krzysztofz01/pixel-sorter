package img

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/disintegration/imaging"
)

// TODO: Fine tune const values in order to make the edge detection more versatile
// TODO: Implement parallel image iteration for NonMaxSuppresion, Hysterysis, NewWhiteNRGBA and ImageGradienPoitns validation
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
func PerformEdgeDetection(i image.Image, performNonMaxSupression bool) (image.Image, error) {
	imgGrayscale := imaging.Grayscale(i)

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

	drawableResult, err := utils.GetDrawableImage(imgHysteresis)
	if err != nil {
		return nil, fmt.Errorf("edge-detection: failed to convert the edge detection image to drawable version: %w", err)
	}

	return drawableResult, nil
}

// Helper function used to create a gradient point matrix by calculating the sobel vertical and horizontal derivatives
func calculateGradientPoints(vertical *image.NRGBA, horizontal *image.NRGBA) ([][]gradientPoint, error) {
	width := vertical.Bounds().Dx()
	height := vertical.Bounds().Dy()

	gp := make([][]gradientPoint, width)
	magnitudeMax := math.Inf(-1)
	for xIndex := 0; xIndex < width; xIndex += 1 {
		gp[xIndex] = make([]gradientPoint, height)
		for yIndex := 0; yIndex < height; yIndex += 1 {
			vColor := float64(utils.ColorToGrayscaleComponent(vertical.NRGBAAt(xIndex, yIndex)))
			hColor := float64(utils.ColorToGrayscaleComponent(horizontal.NRGBAAt(xIndex, yIndex)))

			magnitude := math.Hypot(vColor, hColor)
			if magnitude > magnitudeMax {
				magnitudeMax = magnitude
			}

			direction := math.Atan2(vColor, hColor)
			gp[xIndex][yIndex] = gradientPoint{magnitude: magnitude, direction: direction}
		}
	}

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
	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			magnitude := g[xIndex][yIndex].magnitude
			colorValue := uint8(math.Max(0, math.Min(255, magnitude)))

			outputImage.Set(xIndex, yIndex, color.Gray{Y: colorValue})
		}
	}

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
	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			magnitude := g[xIndex][yIndex].magnitude
			colorValue := uint8(math.Max(0, math.Min(255, magnitude)))

			if xIndex == 0 || xIndex == width-1 || yIndex == 0 || yIndex == height-1 {
				outputImage.Set(xIndex, yIndex, color.Gray{Y: colorValue})
				continue
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
				outputImage.Set(xIndex, yIndex, color.White)
			} else {
				outputImage.Set(xIndex, yIndex, color.Gray{Y: colorValue})
			}
		}
	}

	return outputImage, nil
}

// Helper function used to perform a Double Threshold and Hysteresis operations on the image
func performHysteresis(i *image.NRGBA) *image.NRGBA {
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	maxColorValue := 0
	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			colorValue := utils.ColorToGrayscaleComponent(i.NRGBAAt(xIndex, yIndex))
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
	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			colorValue := float64(utils.ColorToGrayscaleComponent(i.NRGBAAt(xIndex, yIndex)))

			targetColor := colorWeak
			if colorValue >= dtUpperThreshold {
				targetColor = colorStrong
			}

			if colorValue < dtLowerThreshold {
				targetColor = colorZero
			}

			outputImage.Set(xIndex, yIndex, targetColor)
		}
	}

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			if xIndex == 0 || xIndex == width-1 || yIndex == 0 || yIndex == height-1 {
				outputImage.Set(xIndex, yIndex, colorZero)
				continue
			}

			colorValue := uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex, yIndex)))
			if colorValue == doubleThresholdWeak {
				neighbors := [8]uint8{
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex-1, yIndex-1))),
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex-1, yIndex))),
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex-1, yIndex+1))),
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex, yIndex-1))),
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex, yIndex+1))),
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex+1, yIndex-1))),
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex+1, yIndex))),
					uint8(utils.ColorToGrayscaleComponent(outputImage.NRGBAAt(xIndex+1, yIndex+1))),
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
	for xIndex := 0; xIndex < r.Dx(); xIndex += 1 {
		for yIndex := 0; yIndex < r.Dy(); yIndex += 1 {
			img.Set(xIndex, yIndex, color.White)
		}
	}

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
