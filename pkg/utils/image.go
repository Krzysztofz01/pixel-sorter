package utils

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"

	"github.com/Krzysztofz01/pimit"
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

// Create a NRGBA image copy of an image represented as a image.Image.
func ImageToNrgbaImage(i image.Image) *image.NRGBA {
	switch i0 := i.(type) {
	case *image.NRGBA:
		return i0
	case *image.RGBA:
		return RgbaToNrgbaImage(i0)
	default:
		img := image.NewNRGBA(image.Rect(0, 0, i.Bounds().Dx(), i.Bounds().Dy()))
		draw.Draw(img, img.Bounds(), i, i.Bounds().Min, draw.Src)

		return img
	}
}

// Create a RGBA image copy of an image represented as a image.Image.
func ImageToRgbaImage(i image.Image) *image.RGBA {
	switch i0 := i.(type) {
	case *image.RGBA:
		return i0
	case *image.NRGBA:
		return NrgbaToRgbaImage(i0)
	default:
		img := image.NewRGBA(image.Rect(0, 0, i.Bounds().Dx(), i.Bounds().Dy()))
		draw.Draw(img, img.Bounds(), i, i.Bounds().Min, draw.Src)

		return img
	}
}

// Get a copy of the provided image.NRGBA image. This function will panic if the provided image pointer is nil.
func GetImageCopyNrgba(i *image.NRGBA) *image.NRGBA {
	if i == nil {
		panic("image-utils: can not create a copy of a nil nrgba image")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, i.Pix)

	return img
}

// Get a copy of the provided image.RGBA image. This function will panic if the provided image pointer is nil.
func GetImageCopyRgba(i *image.RGBA) *image.RGBA {
	if i == nil {
		panic("image-utils: can not create a copy of a nil rgba image")
	}

	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	copy(img.Pix, i.Pix)

	return img
}

// Get the drawable copy of the provided image. The image is also redrawn to a new image.RGBA struct
//
// Deprecated: GetImgageCopy...
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
//
// Deprecated: InvertImageNrgba
func InvertImage(i image.Image) (draw.Image, error) {
	drawableImage, err := GetDrawableImage(i)
	if err != nil {
		return nil, fmt.Errorf("image-utils: can not get the drawable image version: %w", err)
	}

	pimit.ParallelColumnColorReadWrite(drawableImage, func(c color.Color) color.Color {
		currentColor := ColorToRgba(c)

		return color.RGBA{
			R: 255 - currentColor.R,
			G: 255 - currentColor.G,
			B: 255 - currentColor.B,
			A: currentColor.A,
		}
	})

	return drawableImage, nil
}

// Invert the colors (negative-effect) of the given NRGBA image.
func InvertImageNrgba(i *image.NRGBA) *image.NRGBA {
	if i == nil {
		panic("image-utils: can not invert nil image")
	}

	return pimit.ParallelNrgbaReadWriteNew(i, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return 255 - r, 255 - g, 255 - b, a
	})
}

// Rotate the image by a given angle
//
// Beacuse the dependency internal rotate implementation is using a custom pixel color handling
// solution we need to redraw the result image to ensure that the colors space is RGBA
//
// Deprecated: RotateImageNrgba
func RotateImage(i image.Image, angle int) draw.Image {
	angleNorm := float64(angle) + math.Ceil(-float64(angle)/360.0)*360.0

	rotatedImage := imaging.Rotate(i, angleNorm, color.Transparent)
	rotatedImageWidth := rotatedImage.Bounds().Dx()
	rotatedImageHeight := rotatedImage.Bounds().Dy()

	redrawnImage := image.NewRGBA(image.Rect(0, 0, rotatedImageWidth, rotatedImageHeight))
	draw.Draw(redrawnImage, rotatedImage.Bounds(), rotatedImage, rotatedImage.Bounds().Min, draw.Src)

	return redrawnImage
}

// Rotate the NRGBA image counter-clockwise by a given angle.
func RotateImageNrgba(i *image.NRGBA, angle int) *image.NRGBA {
	angleNorm := float64(angle) + math.Ceil(-float64(angle)/360.0)*360.0
	if angleNorm == 0 {
		return GetImageCopyNrgba(i)
	}

	return imaging.Rotate(i, angleNorm, color.Transparent)
}

// Rotate the NRGBA image counter-clockwise by a given angle and expose a function that can be used
// to rotate the image back to its previous state and remove excess workspace transparency.
func RotateImageWithRevertNrgba(i *image.NRGBA, angle int) (*image.NRGBA, func(*image.NRGBA) *image.NRGBA) {
	if i == nil {
		panic("image-utils: can not perform rotation on a nil image")
	}

	var (
		rotated   *image.NRGBA
		angleNorm float64 = float64(angle) + math.Ceil(-float64(angle)/360.0)*360.0
	)

	if angleNorm == 0 {
		rotated = GetImageCopyNrgba(i)
	} else {
		rotated = imaging.Rotate(i, angleNorm, color.Transparent)
	}

	revertFunc := func(revertImage *image.NRGBA) *image.NRGBA {
		if revertImage.Bounds().Dx() != rotated.Bounds().Dx() || revertImage.Bounds().Dy() != rotated.Bounds().Dy() {
			panic("image-utils: can not revert the image rotation due to invalid bounds")
		}

		if angleNorm == 0 {
			return GetImageCopyNrgba(revertImage)
		}

		r := imaging.Rotate(i, -angleNorm, color.Transparent)
		return trimImageTransparentWorkspaceNrgba(r, i.Bounds())
	}

	return rotated, revertFunc
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
//
// Deprecated: RotateImageWithRevertNrgba
func TrimImageTransparentWorkspace(imageWithWorkspace draw.Image, imageOriginal image.Image) draw.Image {
	xIndexStart := (imageWithWorkspace.Bounds().Dx() - imageOriginal.Bounds().Dx()) / 2
	xIndexLength := imageOriginal.Bounds().Dx()

	yIndexStart := (imageWithWorkspace.Bounds().Dy() - imageOriginal.Bounds().Dy()) / 2
	yIndexLength := imageOriginal.Bounds().Dy()

	if xIndexLength == 0 && yIndexLength == 0 {
		return imageWithWorkspace
	}

	tImg := image.NewRGBA(image.Rect(0, 0, xIndexLength, yIndexLength))
	pimit.ParallelColumnReadWrite(tImg, func(xIndex, yIndex int, _ color.Color) color.Color {
		xOffset := xIndex + xIndexStart
		yOffset := yIndex + yIndexStart

		return imageWithWorkspace.At(xOffset, yOffset)
	})

	return tImg
}

func trimImageTransparentWorkspaceNrgba(withWorkspace *image.NRGBA, original image.Rectangle) *image.NRGBA {
	var (
		xIndexLength int = original.Dx()
		yIndexLength int = original.Dy()
		xIndexStart  int = (withWorkspace.Bounds().Dx() - xIndexLength) / 2
		yIndexStart  int = (withWorkspace.Bounds().Dy() - yIndexLength) / 2
	)

	if xIndexStart == 0 && yIndexStart == 0 {
		return GetImageCopyNrgba(withWorkspace)
	}

	img := image.NewNRGBA(image.Rect(0, 0, xIndexLength, yIndexLength))
	pimit.ParallelNrgbaReadWrite(img, func(x, y int, _, _, _, _ uint8) (uint8, uint8, uint8, uint8) {
		xOffset := x + xIndexStart
		yOffset := y + yIndexStart
		index := 4 * (yOffset*xIndexLength + xOffset)

		return withWorkspace.Pix[index+0], withWorkspace.Pix[index+1], withWorkspace.Pix[index+2], withWorkspace.Pix[index+3]
	})

	return img
}

// Function used to scale the image down according to given percentage parameter (Value from 0.0 to 1.0)
//
// Deprecated: ScaleImageNrgba
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

// Function used to scale a NRGBA image down according to given percentage parameter (Value from 0.0 to 1.0).
func ScaleImageNrgba(i *image.NRGBA, percentage float64) (*image.NRGBA, error) {
	if percentage < 0.0 || percentage > 1.0 {
		return nil, errors.New("image-utils: invalid downscale percentage specified")
	}

	if percentage == 1.0 {
		return GetImageCopyNrgba(i), nil
	}

	scaledWidth := int(float64(i.Bounds().Dx()) * percentage)
	scaledHeight := int(float64(i.Bounds().Dy()) * percentage)

	// NOTE: The usage of the "Lanczos" algorithm will produce images with better quality, for perfomance use "Box"
	return imaging.Resize(i, scaledWidth, scaledHeight, imaging.Lanczos), nil
}

// Blend two images using a given blending mode into a new image
//
// Deprecated: BlendImagesNrgba
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
	pimit.ParallelColumnReadWrite(resultImage, func(xIndex, yIndex int, _ color.Color) color.Color {
		aColor := ColorToRgba(aRgba.At(xIndex, yIndex))
		bColor := ColorToRgba(bRgba.At(xIndex, yIndex))

		return BlendRGBA(aColor, bColor, mode)
	})

	return resultImage, nil
}

// Blend two NRGBA images using a given blending mode into a new image.
// TODO: Reimplement using NRGBA values and internal buffers of src params
func BlendImagesNrgba(a, b *image.NRGBA, mode BlendingMode) (*image.NRGBA, error) {
	if a == nil || b == nil {
		panic("image-utils: can not perform blending if one of the images is nil")
	}

	if a.Bounds().Dx() != b.Bounds().Dx() {
		return nil, errors.New("image-utils: the provided images have different width")
	}

	if a.Bounds().Dy() != b.Bounds().Dy() {
		return nil, errors.New("image-utils: the provided images have different height")
	}

	img := image.NewNRGBA(image.Rect(0, 0, a.Bounds().Dx(), a.Bounds().Dy()))
	pimit.ParallelReadWrite(img, func(x, y int, c color.Color) color.Color {
		aColor := ColorToRgba(a.NRGBAAt(x, y))
		bColor := ColorToRgba(b.NRGBAAt(x, y))

		return BlendRGBA(aColor, bColor, mode)
	})

	return img, nil
}

// Create a RGBA color-space copy of an image represented in the NRGBA color-space.
func NrgbaToRgbaImage(i *image.NRGBA) *image.RGBA {
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	pimit.ParallelNrgbaRead(i, func(x, y int, r, g, b, a uint8) {
		index := 4 * (y*width + x)

		a32 := uint32(a) * 0x101
		r32 := uint32(r) * a32 / 0xff
		g32 := uint32(g) * a32 / 0xff
		b32 := uint32(b) * a32 / 0xff

		img.Pix[index+0] = uint8(r32 >> 8)
		img.Pix[index+1] = uint8(g32 >> 8)
		img.Pix[index+2] = uint8(b32 >> 8)
		img.Pix[index+3] = uint8(a32 >> 8)
	})

	return img
}

// Create a NRGBA color-space copy of an image represented in the RGBA color-space.
// TODO: Reimplement this function using the pimit parallelization
func RgbaToNrgbaImage(i *image.RGBA) *image.NRGBA {
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), i, i.Bounds().Min, draw.Src)
	return img
}
