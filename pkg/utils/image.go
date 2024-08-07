package utils

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/Krzysztofz01/imaging"
	"github.com/Krzysztofz01/pimit"
)

// Convert a image represented by image.Image to a NRGBA image. If the underlying type is already NRGBA a reference to the
// original input will be returned, otherwis a copy will be created.
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

// Convert a image represented by image.Image to a RGBA image. If the underlying type is already RGBA a reference to the
// original input will be returned, otherwis a copy will be created.
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

// Invert the colors (negative-effect) of the given NRGBA image.
func InvertImageNrgba(i *image.NRGBA) *image.NRGBA {
	if i == nil {
		panic("image-utils: can not invert nil image")
	}

	return pimit.ParallelNrgbaReadWriteNew(i, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		return 255 - r, 255 - g, 255 - b, a
	})
}

// Rotate the NRGBA image counter-clockwise by a given angle.
func RotateImageNrgba(i *image.NRGBA, angle int) *image.NRGBA {
	if i == nil {
		panic("image-utils: can not perform rotation on a nil image")
	}

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

	var (
		originalImageBounds image.Rectangle = i.Bounds()
		rotatedImageBounds  image.Rectangle = rotated.Bounds()
	)

	revertFunc := func(revertImage *image.NRGBA) *image.NRGBA {
		if revertImage == nil {
			panic("image-utils: can not perform rotation revert on a nil image")
		}

		if revertImage.Bounds().Dx() != rotatedImageBounds.Dx() || revertImage.Bounds().Dy() != rotatedImageBounds.Dy() {
			panic("image-utils: can not revert the image rotation due to invalid bounds")
		}

		if angleNorm == 0 {
			return GetImageCopyNrgba(revertImage)
		}

		r := imaging.Rotate(revertImage, -angleNorm, color.Transparent)
		return trimImageTransparentWorkspaceNrgba(r, originalImageBounds)
	}

	return rotated, revertFunc
}

// TODO: Add missing docs and unit tests
func TrimImageTransparentWorkspaceNrgba(withWorkspace *image.NRGBA, original image.Rectangle) *image.NRGBA {
	return trimImageTransparentWorkspaceNrgba(withWorkspace, original)
}

func trimImageTransparentWorkspaceNrgba(withWorkspace *image.NRGBA, original image.Rectangle) *image.NRGBA {
	var (
		originalWidth   int = original.Dx()
		originalHeight  int = original.Dy()
		workspaceWidth  int = withWorkspace.Bounds().Dx()
		workspaceHeight int = withWorkspace.Bounds().Dy()
		xIndexStart     int = (workspaceWidth - originalWidth) / 2
		yIndexStart     int = (workspaceHeight - originalHeight) / 2
	)

	img := image.NewNRGBA(image.Rect(0, 0, originalWidth, originalHeight))
	pimit.ParallelNrgbaReadWrite(img, func(x, y int, _, _, _, _ uint8) (uint8, uint8, uint8, uint8) {
		offset := 4 * (workspaceWidth*(yIndexStart+y) + xIndexStart + x)

		return withWorkspace.Pix[offset+0], withWorkspace.Pix[offset+1], withWorkspace.Pix[offset+2], withWorkspace.Pix[offset+3]
	})

	return img
}

// Function used to scale a NRGBA image down according to given percentage parameter (Value from 0.0 to 1.0).
func ScaleImageNrgba(i *image.NRGBA, percentage float64) (*image.NRGBA, error) {
	if percentage <= 0.0 || percentage > 1.0 {
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

// Blend two NRGBA images using a given blending mode into a new image.
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
	pimit.ParallelNrgbaReadWrite(img, func(x, y int, _, _, _, _ uint8) (uint8, uint8, uint8, uint8) {
		color := BlendNrgba(a.NRGBAAt(x, y), b.NRGBAAt(x, y), mode)
		return color.R, color.G, color.B, color.A
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

// Create a grayscale version of the provided NRGBA color-space image.
func GrayscaleNrgba(i *image.NRGBA) *image.NRGBA {
	return pimit.ParallelNrgbaReadWriteNew(i, func(x, y int, r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		gray := uint8((float64(r) * 0.299) + (float64(g) * 0.587) + (float64(b) * 0.114))
		return gray, gray, gray, a
	})
}
