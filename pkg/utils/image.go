package utils

import (
	"errors"
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

// Get the drawable version of the provided image. The image is also redrawn to a new image.RGBA struct
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

// Rotate the image by a given angle
func RotateImage(image draw.Image, angle int) draw.Image {
	angleNorm := float64(angle) + math.Ceil(-float64(angle)/360.0)*360.0

	return imaging.Rotate(image, angleNorm, color.Transparent)
}

// Remove all excess transparent workspace created during the rotation process from the image.
// This function is detection the top left corner on the rotatet image and than is taking cropping
// the image by width and height given by the original image.
// Time complexity: O(2n)
// NOTE: This function can by changed to detect the content without the imageOriginal
func TrimImageTransparentWorkspace(imageWithWorkspace draw.Image, imageOriginal draw.Image) draw.Image {
	wImgWidth := imageWithWorkspace.Bounds().Dx()
	wImgHeight := imageWithWorkspace.Bounds().Dy()

	tR, tG, tB, _ := imageOriginal.At(0, 0).RGBA()
	xOffset := 0
	yOffset := 0

	for yIndex := 0; yIndex < wImgHeight; yIndex += 1 {
		for xIndex := 0; xIndex < wImgWidth; xIndex += 1 {
			cR, cG, cB, _ := imageWithWorkspace.At(xIndex, yIndex).RGBA()
			if tR == cR && tG == cG && tB == cB {
				yOffset = yIndex
				xOffset = xIndex
				break
			}
		}
	}

	oImgWidth := imageOriginal.Bounds().Dx()
	oImgHeight := imageOriginal.Bounds().Dy()
	tImg := image.NewRGBA(image.Rect(0, 0, oImgWidth, oImgHeight))

	for yIndex := 0; yIndex < oImgHeight; yIndex += 1 {
		for xIndex := 0; xIndex < oImgWidth; xIndex += 1 {
			color := imageWithWorkspace.At(xIndex+xOffset, yIndex+yOffset)
			tImg.Set(xIndex, yIndex, color)
		}
	}

	return tImg
}
