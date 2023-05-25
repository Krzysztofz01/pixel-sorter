package utils

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"

	"github.com/disintegration/imaging"
)

// Get the image column as a color.Color interface implementation slice specified by the x image index. The retrieval
// process is run parallel in several goroutines.
func GetImageColumn(image image.Image, xIndex int) ([]color.Color, error) {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	if xIndex >= width {
		return nil, errors.New("image-utils: column index is not in range of the target image width")
	}

	iteratorCount := height / 400
	if iteratorCount < 1 {
		iteratorCount = 1
	}

	lengthBase := height / iteratorCount
	lengthRemnant := height % iteratorCount

	wg := sync.WaitGroup{}
	wg.Add(iteratorCount)

	column := make([]color.Color, height)
	for offsetFactor := 0; offsetFactor < iteratorCount; offsetFactor += 1 {
		targetOffset := offsetFactor * lengthBase
		targetLength := lengthBase
		if offsetFactor+1 == iteratorCount {
			targetLength += lengthRemnant
		}

		go func(offset, length int) {
			defer wg.Done()

			for iterationOffset := 0; iterationOffset < length; iterationOffset += 1 {
				currentOffset := offset + iterationOffset

				color := image.At(xIndex, currentOffset)
				column[currentOffset] = color
			}
		}(targetOffset, targetLength)
	}

	wg.Wait()
	return column, nil
}

// Get the image row as a color.Color interface implementation slice specified by the y image index. The retrieval
// process is run parallel in several goroutines.
func GetImageRow(image image.Image, yIndex int) ([]color.Color, error) {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	if yIndex >= height {
		return nil, errors.New("image-utils: row index is not in range of the target image height")
	}

	iteratorCount := width / 400
	if iteratorCount < 1 {
		iteratorCount = 1
	}

	lengthBase := width / iteratorCount
	lengthRemnant := width % iteratorCount

	wg := sync.WaitGroup{}
	wg.Add(iteratorCount)

	row := make([]color.Color, width)
	for offsetFactor := 0; offsetFactor < iteratorCount; offsetFactor += 1 {
		targetOffset := offsetFactor * lengthBase
		targetLength := lengthBase
		if offsetFactor+1 == iteratorCount {
			targetLength += lengthRemnant
		}

		go func(offset, length int) {
			defer wg.Done()

			for iterationOffset := 0; iterationOffset < length; iterationOffset += 1 {
				currentOffset := offset + iterationOffset

				color := image.At(currentOffset, yIndex)
				row[currentOffset] = color
			}
		}(targetOffset, targetLength)
	}

	wg.Wait()
	return row, nil
}

// Set the colors of the column of the given image specified by the xIndex, according to the slice containing implementations
// of the color.Color interface. The changes made will be applied to the given draw.Image reference. The write process is
// run parallel in several goroutines.
func SetImageColumn(image draw.Image, column []color.Color, xIndex int) error {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	if xIndex >= width {
		return errors.New("image-utils: column index is not in range of the target image width")
	}

	if len(column) != height {
		return errors.New("image-utils: the image height and the provided column lengths are not matching")
	}

	iteratorCount := height / 400
	if iteratorCount < 1 {
		iteratorCount = 1
	}

	lengthBase := height / iteratorCount
	lengthRemnant := height % iteratorCount

	wg := sync.WaitGroup{}
	wg.Add(iteratorCount)

	for offsetFactor := 0; offsetFactor < iteratorCount; offsetFactor += 1 {
		targetOffset := offsetFactor * lengthBase
		targetLength := lengthBase
		if offsetFactor+1 == iteratorCount {
			targetLength += lengthRemnant
		}

		go func(offset, length int) {
			defer wg.Done()

			for iterationOffset := 0; iterationOffset < length; iterationOffset += 1 {
				currentOffset := offset + iterationOffset

				image.Set(xIndex, currentOffset, column[currentOffset])
			}
		}(targetOffset, targetLength)
	}

	wg.Wait()
	return nil
}

// Set the colors of the row of the given image specified by the yIndex, according to the slice containing implementations
// of the color.Color interface. The changes made will be applied to the given draw.Image reference. The write process is
// run parallel in several goroutines.
func SetImageRow(image draw.Image, row []color.Color, yIndex int) error {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	if yIndex >= height {
		return errors.New("image-utils: row index is not in range of the target image height")
	}

	if len(row) != width {
		return errors.New("image-utils: the image widht and the provided row lengths are not matching")
	}

	iteratorCount := width / 400
	if iteratorCount < 1 {
		iteratorCount = 1
	}

	lengthBase := width / iteratorCount
	lengthRemnant := width % iteratorCount

	wg := sync.WaitGroup{}
	wg.Add(iteratorCount)

	for offsetFactor := 0; offsetFactor < iteratorCount; offsetFactor += 1 {
		targetOffset := offsetFactor * lengthBase
		targetLength := lengthBase
		if offsetFactor+1 == iteratorCount {
			targetLength += lengthRemnant
		}

		go func(offset, length int) {
			defer wg.Done()

			for iterationOffset := 0; iterationOffset < length; iterationOffset += 1 {
				currentOffset := offset + iterationOffset

				image.Set(currentOffset, yIndex, row[currentOffset])
			}
		}(targetOffset, targetLength)
	}

	wg.Wait()
	return nil
}

// Get the drawable copy of the provided image. The image is also redrawn to a new image.RGBA struct
func GetDrawableImage(i image.Image) (draw.Image, error) {
	p1 := image.Point{0, 0}
	p2 := image.Point{i.Bounds().Dx(), i.Bounds().Dy()}

	rgbaImage := image.NewRGBA(image.Rectangle{p1, p2})
	draw.Draw(rgbaImage, rgbaImage.Bounds(), i, i.Bounds().Min, draw.Src)

	// NOTE: As we always convert the image to the RGBA image which IS implementing the draw.Image
	// this check becomes redundant. We may keep it here in case we drop the idea of redrawing the
	// image for performance purposes.
	//
	// if _, ok := rgbaImage.(*draw.Image); !ok {
	// 	return nil, errors.New("image-utils: unable to create a drawable image")
	// }

	return rgbaImage, nil
}

// Invert the colors (negative-effect) of the image
func InvertImage(i image.Image) (draw.Image, error) {
	drawableImage, err := GetDrawableImage(i)
	if err != nil {
		return nil, fmt.Errorf("image-utils: can not get the drawable image version: %w", err)
	}

	width := drawableImage.Bounds().Dx()
	height := drawableImage.Bounds().Dy()

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			currentColor := ColorToRgba(drawableImage.At(xIndex, yIndex))

			invertedColor := color.RGBA{
				R: 255 - currentColor.R,
				G: 255 - currentColor.G,
				B: 255 - currentColor.B,
				A: currentColor.A,
			}

			drawableImage.Set(xIndex, yIndex, invertedColor)
		}
	}

	return drawableImage, nil
}

// Rotate the image by a given angle
//
// Beacuse the dependency internal rotate implementation is using a custom pixel color handling
// solution we need to redraw the result image to ensure that the colors space is RGBA
func RotateImage(i draw.Image, angle int) draw.Image {
	angleNorm := float64(angle) + math.Ceil(-float64(angle)/360.0)*360.0

	rotatedImage := imaging.Rotate(i, angleNorm, color.Transparent)
	rotatedImageWidth := rotatedImage.Bounds().Dx()
	rotatedImageHeight := rotatedImage.Bounds().Dy()

	redrawnImage := image.NewRGBA(image.Rect(0, 0, rotatedImageWidth, rotatedImageHeight))
	draw.Draw(redrawnImage, rotatedImage.Bounds(), rotatedImage, rotatedImage.Bounds().Min, draw.Src)

	return redrawnImage
}

// FIXME: Fix borders after interpolation and implement full transparency support in order to preserve PNG data.
// FIXME: Better handling for situations where there is not requirement to trim the workspace. Currently in such
// situation we are returning the original imageWithWorkspace draw.Image instead of creating a copy.
//
// Remove all excess transparent workspace created during the rotation process from the image.
// The function is calculating the image located in the middle of the workspace and is cropping it.
// Due to rotation interpolation the borders of the image are not matching the original image,
// and the whole operation can cause some problems realted to transparent colors.
//
// Time complexity: O(n)
func TrimImageTransparentWorkspace(imageWithWorkspace draw.Image, imageOriginal image.Image) draw.Image {
	xIndexStart := (imageWithWorkspace.Bounds().Dx() - imageOriginal.Bounds().Dx()) / 2
	xIndexLength := imageOriginal.Bounds().Dx()

	yIndexStart := (imageWithWorkspace.Bounds().Dy() - imageOriginal.Bounds().Dy()) / 2
	yIndexLength := imageOriginal.Bounds().Dy()

	if xIndexLength == 0 && yIndexLength == 0 {
		return imageWithWorkspace
	}

	tImg := image.NewRGBA(image.Rect(0, 0, xIndexLength, yIndexLength))

	for xIndex := 0; xIndex < xIndexLength; xIndex += 1 {
		for yIndex := 0; yIndex < yIndexLength; yIndex += 1 {
			xOffset := xIndex + xIndexStart
			yOffset := yIndex + yIndexStart

			color := imageWithWorkspace.At(xOffset, yOffset)
			tImg.Set(xIndex, yIndex, color)
		}
	}

	return tImg
}

// Function used to scale the image down according to given percentage parameter (Value from 0.0 to 1.0)
func ScaleImage(i image.Image, percentage float64) (draw.Image, error) {
	if percentage < 0.0 || percentage > 1.0 {
		return nil, errors.New("image-utils: invalid downscale percentage specified")
	}

	scaledWidth := float64(i.Bounds().Dx()) * percentage
	scaledHeight := float64(i.Bounds().Dy()) * percentage

	// NOTE: The usage of the "Lanczos" algorithm will produce images with better quality, for perfomance use "Box"
	scaledImage := imaging.Resize(i, int(scaledWidth), int(scaledHeight), imaging.Lanczos)

	scaledDrawableImage, err := GetDrawableImage(scaledImage)
	if err != nil {
		return nil, fmt.Errorf("image-utils: can not get the scaled drawable image version: %w", err)
	}

	return scaledDrawableImage, nil
}

// Blend two images using a given blending mode into a new image
// TODO: Unit test implementation
func BlendImages(a, b image.Image, mode BlendingMode) (draw.Image, error) {
	aWidth := a.Bounds().Dx()
	aHeight := a.Bounds().Dy()

	bWidth := b.Bounds().Dx()
	bHeight := b.Bounds().Dy()

	if aWidth != bWidth {
		return nil, errors.New("image-utils: the provided images have a different width")
	}

	if aHeight != bHeight {
		return nil, errors.New("image-utils: the provided images have a different height")
	}

	aRgba, err := GetDrawableImage(a)
	if err != nil {
		return nil, fmt.Errorf("image-utils: can not get the drawable image version of a: %w", err)
	}

	bRgba, err := GetDrawableImage(b)
	if err != nil {
		return nil, fmt.Errorf("image-utils: can not get the drawable image version of b: %w", err)
	}

	resultImage := image.NewRGBA(image.Rect(0, 0, aWidth, aHeight))

	for xIndex := 0; xIndex < aWidth; xIndex += 1 {
		for yIndex := 0; yIndex < aHeight; yIndex += 1 {
			aColor := ColorToRgba(aRgba.At(xIndex, yIndex))
			bColor := ColorToRgba(bRgba.At(xIndex, yIndex))

			resultImage.Set(xIndex, yIndex, BlendRGBA(aColor, bColor, mode))
		}
	}

	return resultImage, nil
}
