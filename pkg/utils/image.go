package utils

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/disintegration/imaging"
)

// Get the image column as a Color interface slice specified by the x image index
func GetImageColumn(image image.Image, xIndex int) ([]color.Color, error) {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	if xIndex >= width {
		return nil, errors.New("image-utils: column index is not in range of the target image width")
	}

	column := make([]color.Color, height)
	for yIndex := 0; yIndex < height; yIndex += 1 {
		color := image.At(xIndex, yIndex)
		column[yIndex] = color
	}

	return column, nil
}

// Get the image row as a Color interface slice specified by the y image index
func GetImageRow(image image.Image, yIndex int) ([]color.Color, error) {
	width := image.Bounds().Dx()
	height := image.Bounds().Dy()

	if yIndex >= height {
		return nil, errors.New("image-utils: row index is not in range of the target image height")
	}

	row := make([]color.Color, width)
	for xIndex := 0; xIndex < width; xIndex += 1 {
		color := image.At(xIndex, yIndex)
		row[xIndex] = color
	}

	return row, nil
}

// Append a Color interface slice representing a column at a given x index of the given image
func SetImageColumn(image *draw.Image, column []color.Color, xIndex int) error {
	width := (*image).Bounds().Dx()
	height := (*image).Bounds().Dy()

	if xIndex >= width {
		return errors.New("image-utils: column index is not in range of the target image width")
	}

	if len(column) != height {
		return errors.New("image-utils: the image height and the provided column lengths are not matching")
	}

	for yIndex := 0; yIndex < height; yIndex += 1 {
		color := column[yIndex]
		(*image).Set(xIndex, yIndex, color)
	}

	return nil
}

// Append a Color interface slice representing a row at a given y index of the given image
func SetImageRow(image *draw.Image, row []color.Color, yIndex int) error {
	width := (*image).Bounds().Dx()
	height := (*image).Bounds().Dy()

	if yIndex >= height {
		return errors.New("image-utils: row index is not in range of the target image height")
	}

	if len(row) != width {
		return errors.New("image-utils: the image widht and the provided row lengths are not matching")
	}

	for xIndex := 0; xIndex < width; xIndex += 1 {
		color := row[xIndex]
		(*image).Set(xIndex, yIndex, color)
	}

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

func InvertImage(i image.Image) (draw.Image, error) {
	drawableImage, err := GetDrawableImage(i)
	if err != nil {
		return nil, fmt.Errorf("image-utils: can not get the drawable image version: %w", err)
	}

	width := drawableImage.Bounds().Dx()
	height := drawableImage.Bounds().Dy()

	for xIndex := 0; xIndex < width; xIndex += 1 {
		for yIndex := 0; yIndex < height; yIndex += 1 {
			currentColor, err := ColorToRgba(drawableImage.At(xIndex, yIndex))
			if err != nil {
				return nil, fmt.Errorf("image-utils: can not access the color as RGBA: %w", err)
			}

			drawableImage.Set(xIndex, yIndex, color.RGBA{
				R: 255 - currentColor.R,
				G: 255 - currentColor.G,
				B: 255 - currentColor.B,
				A: currentColor.A,
			})
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
