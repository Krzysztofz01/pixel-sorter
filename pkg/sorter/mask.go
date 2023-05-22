package sorter

import (
	"errors"
	"image"
	"image/color"

	"github.com/Krzysztofz01/pimit"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

// A structure representing a wrapper over a image that works as a mask
type Mask struct {
	maskImage image.Image
	isEmpty   bool
}

// Crate a new mask instance from a given image without target image restrictions
func CreateMask(mImg image.Image) (*Mask, error) {
	return CreateImageMask(mImg, mImg.Bounds())
}

// Create a new mask instance from a given image and bounds of the image to be masked
func CreateImageMask(mImg image.Image, targetImageBounds image.Rectangle) (*Mask, error) {
	if mImg.Bounds().Dx() != targetImageBounds.Dx() || mImg.Bounds().Dy() != targetImageBounds.Dy() {
		return nil, errors.New("sorter: mask image and target image sizes are not matching")
	}

	err := pimit.ParallelColumnColorReadE(mImg, func(c color.Color) error {
		currentColor := utils.ColorToRgba(c)
		_, s, l := utils.RgbaToHsl(currentColor)

		if s != 0.0 || (l != 0.0 && l != 1.0) {
			return errors.New("sorter: the mask contains a invalid color")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	mask := new(Mask)
	mask.maskImage = mImg
	mask.isEmpty = false
	return mask, nil
}

// Create a new mask instance representing a empty mask
func CreateEmptyMask() *Mask {
	mask := new(Mask)
	mask.maskImage = nil
	mask.isEmpty = true
	return mask
}

// Perform a mask lookup to check if the mask is masking at the given location
func (mask *Mask) IsMasked(xIndex, yIndex int) (bool, error) {
	// TODO: If the mask is empty, the size is not validate. We should do something about it in the future, if we want to perform size validation in mask factory func
	if mask.isEmpty {
		return false, nil
	}

	xLength := mask.maskImage.Bounds().Dx()
	yLength := mask.maskImage.Bounds().Dy()
	if xIndex >= xLength || yIndex >= yLength {
		return false, errors.New("sorter: the mask lookup is out of the mask bounds")
	}

	color := utils.ColorToRgba(mask.maskImage.At(xIndex, yIndex))
	_, _, l := utils.RgbaToHsl(color)

	if l == 0.0 {
		return true, nil
	} else if l == 1.0 {
		return false, nil
	} else {
		return false, errors.New("sorter: the mask lookup found a invalid value and may be corrupted")
	}
}
