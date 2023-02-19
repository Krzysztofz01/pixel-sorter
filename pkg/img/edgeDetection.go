package img

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/disintegration/imaging"
)

// TODO: Fine tune this values
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

type gradientPoint struct {
	magnitude float64
	direction float64
}

func DetectEdgesCanny(i image.Image) (image.Image, error) {
	imgGrayscale := imaging.Grayscale(i)
	// _ = utils.StoreImageToFile("1-grayscale.png", "png", imgGrayscale)

	imgSmoothed := imaging.Blur(imgGrayscale, blurSigmaParam)
	// _ = utils.StoreImageToFile("2-smoothed.png", "png", imgSmoothed)

	imgVerticalConv := imaging.Convolve3x3(imgSmoothed, sobelMatrixVertical, nil)
	// _ = utils.StoreImageToFile("3-vertical-conv.png", "png", imgVerticalConv)

	imgHorizontalConv := imaging.Convolve3x3(imgSmoothed, sobelMatrixHorizontal, nil)
	// _ = utils.StoreImageToFile("3-horizontal-conv.png", "png", imgHorizontalConv)

	gradientPoints, err := calculateGradientPoints(imgVerticalConv, imgHorizontalConv)
	if err != nil {
		return nil, fmt.Errorf("edge-detection: failed to calculate the gradients and slopes: %w", err)
	}
	// dumpGradient(imgVerticalConv, gradientPoints)

	imgNonMaxSuppresion := performNonMaxSuppresion(imgSmoothed, gradientPoints)
	// _ = utils.StoreImageToFile("4-nms.png", "png", imgNonMaxSuppresion)

	imgHysteresis := performHysteresis(imgNonMaxSuppresion, gradientPoints)
	// _ = utils.StoreImageToFile("5-hysteresis.png", "png", imgHysteresis)

	drawableResult, err := utils.GetDrawableImage(imgHysteresis)
	if err != nil {
		return nil, fmt.Errorf("edge-detection: failed to convert the edge detection image to drawable version: %w", err)
	}

	// _ = utils.StoreImageToFile("6-finished.png", "png", drawableResult)
	return drawableResult, nil

}

func calculateGradientPoints(vertical *image.NRGBA, horizontal *image.NRGBA) ([][]gradientPoint, error) {
	width := vertical.Bounds().Dx()
	height := vertical.Bounds().Dy()

	gp := make([][]gradientPoint, width)
	magnitudeMax := math.Inf(-1)
	for xIndex := 0; xIndex < width; xIndex += 1 {
		gp[xIndex] = make([]gradientPoint, height)
		for yIndex := 0; yIndex < height; yIndex += 1 {
			vColor := float64(utils.NrgbaToGrayscaleComponent(vertical.NRGBAAt(xIndex, yIndex)))
			hColor := float64(utils.NrgbaToGrayscaleComponent(horizontal.NRGBAAt(xIndex, yIndex)))

			magnitude := math.Hypot(vColor, hColor)
			if magnitude > magnitudeMax {
				magnitudeMax = magnitude
			}

			direction := math.Atan2(vColor, hColor)
			gp[xIndex][yIndex] = gradientPoint{magnitude: magnitude, direction: direction}
		}
	}

	// for xIndex := 0; xIndex < width; xIndex += 1 {
	// 	for yIndex := 0; yIndex < height; yIndex += 1 {
	// 		// FIXME: Zero division?
	// 		magnitude := gp[xIndex][yIndex].magnitude
	// 		gp[xIndex][yIndex].magnitude = magnitude / magnitudeMax * 255.0
	// 	}
	// }

	return gp, nil
}

// TODO: Gradient slice and image size validation
// TODO: Perform this operation parallel
func performNonMaxSuppresion(i *image.NRGBA, gradients [][]gradientPoint) *image.NRGBA {
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	outputImage := newWhiteNRGBA(image.Rect(0, 0, width, height))
	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			magnitude := gradients[xIndex][yIndex].magnitude
			colorValue := uint8(math.Max(0, math.Min(255, magnitude)))

			if xIndex == 0 || xIndex == width-1 || yIndex == 0 || yIndex == height-1 {
				outputImage.Set(xIndex, yIndex, color.Gray{Y: colorValue})
				continue
			}

			magnitudeIsMax := false
			direction := gradients[xIndex][yIndex].direction * 180.0 / math.Pi
			if direction < 0 {
				direction += 180.0
			}

			if valueBetween(direction, 0, 22.5) || valueBetween(direction, 157.5, 180) {
				if magnitude > gradients[xIndex][yIndex+1].magnitude && magnitude > gradients[xIndex][yIndex-1].magnitude {
					magnitudeIsMax = true
				}
			}

			if valueBetween(direction, 22.5, 67.5) {
				if magnitude > gradients[xIndex+1][yIndex-1].magnitude && magnitude > gradients[xIndex-1][yIndex+1].magnitude {
					magnitudeIsMax = true
				}
			}

			if valueBetween(direction, 67.5, 112.5) {
				if magnitude > gradients[xIndex+1][yIndex].magnitude && magnitude > gradients[xIndex-1][yIndex].magnitude {
					magnitudeIsMax = true
				}
			}

			if valueBetween(direction, 112.5, 157.5) {
				if magnitude > gradients[xIndex-1][yIndex-1].magnitude && magnitude > gradients[xIndex+1][yIndex+1].magnitude {
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

	return outputImage
}

// TODO: Gradient slice and image size validation
// TODO: Perform this operation parallel
func performHysteresis(i *image.NRGBA, gradients [][]gradientPoint) *image.NRGBA {
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
	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			colorValue := float64(utils.NrgbaToGrayscaleComponent(i.NRGBAAt(xIndex, yIndex)))

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

	// utils.StoreImageToFile("44-doublethreshold", "png", outputImage)

	for xIndex := 0; xIndex < width; xIndex += 1 {
		if xIndex == 0 || xIndex == width-1 {
			continue
		}

		for yIndex := 0; yIndex < height; yIndex += 1 {
			if yIndex == 0 || yIndex == height-1 {
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

func valueBetween(value, min, max float64) bool {
	if min > max {
		panic("edge-detection: invalid value between operation boundaries")
	}

	return value > min && value < max
}

func newWhiteNRGBA(r image.Rectangle) *image.NRGBA {
	img := image.NewNRGBA(r)
	for xIndex := 0; xIndex < r.Dx(); xIndex += 1 {
		for yIndex := 0; yIndex < r.Dy(); yIndex += 1 {
			img.Set(xIndex, yIndex, color.White)
		}
	}

	return img
}

// func dumpGradient(i image.Image, g [][]gradientPoint) {
// 	img := image.NewGray(i.Bounds())
// 	for x := 0; x < i.Bounds().Dx(); x += 1 {
// 		for y := 0; y < i.Bounds().Dy(); y += 1 {
// 			mag := math.Min(255, math.Max(0, g[x][y].magnitude))
// 			img.SetGray(x, y, color.Gray{Y: uint8(mag)})
// 		}
// 	}

// 	utils.StoreImageToFile("33-gradient-map", "png", img)
// }
