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
	maskImage           image.Image
	maskImageTranslated image.Image
	isEmpty             bool
	isTranslated        bool
}

// Crate a new mask instance from a given image without target image restrictions
func CreateMask(mImg image.Image) (*Mask, error) {
	if mImg == nil {
		return nil, errors.New("sorter: the provided mask image reference is nil")
	}

	return CreateImageMask(mImg, mImg.Bounds(), 0)
}

// Create a new mask instance from a given image and bounds of the image to be masked. The trnslateAngle parameter
// indicates whether the lookup mask should be interpreted from a given angle.
func CreateImageMask(mImg image.Image, targetImageBounds image.Rectangle, translateAngle int) (*Mask, error) {
	if mImg == nil {
		return nil, errors.New("sorter: the provided mask image reference is nil")
	}

	if mImg.Bounds().Dx() != targetImageBounds.Dx() || mImg.Bounds().Dy() != targetImageBounds.Dy() {
		return nil, errors.New("sorter: mask image and target image sizes are not matching")
	}

	err := pimit.ParallelColumnColorReadE(mImg, func(c color.Color) error {
		// TODO: Refactor to use color.NRGBA
		_, s, l, _ := utils.RgbaToHsla(utils.ColorToRgba(c))

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

	if translateAngle != 0 {
		// TODO: Refactor to use image.NRGBA
		mask.maskImageTranslated = utils.RotateImageNrgba(utils.ImageToNrgbaImage(mImg), translateAngle)
		mask.isTranslated = true
	} else {
		mask.maskImageTranslated = nil
		mask.isTranslated = false
	}

	return mask, nil
}

// Create a new mask instance representing a empty mask
func CreateEmptyMask() *Mask {
	mask := new(Mask)
	mask.maskImage = nil
	mask.maskImageTranslated = nil
	mask.isEmpty = true
	mask.isTranslated = false
	return mask
}

// TODO: If the mask is empty the size is not validated.
// Perform a mask lookup using one-dimentional index to check if the mask is masking at the given location. If the
// mask is empty it will always return false.
func (mask *Mask) IsMaskedByIndex(index int) (bool, error) {
	if mask.isEmpty {
		return false, nil
	}

	width := mask.maskImage.Bounds().Dx()
	if mask.isTranslated {
		width = mask.maskImageTranslated.Bounds().Dx()
	}

	if index < 0 {
		return false, errors.New("sorter: the mask index lookup is out of bounds")
	}

	x := index % width
	y := index / width

	return mask.IsMasked(x, y)
}

// TODO: If the mask is empty the size is not validated.
// Perform a mask lookup to check if the mask is masking at the given location. If the
// mask is empty it will always return false.
func (mask *Mask) IsMasked(xIndex, yIndex int) (bool, error) {
	if mask.isEmpty {
		return false, nil
	}

	var color color.Color
	if mask.isTranslated {
		xLength := mask.maskImageTranslated.Bounds().Dx()
		yLength := mask.maskImageTranslated.Bounds().Dy()
		if xIndex >= xLength || yIndex >= yLength {
			return false, errors.New("sorter: the mask lookup is out of the translated mask bounds")
		}

		color = utils.ColorToRgba(mask.maskImageTranslated.At(xIndex, yIndex))
	} else {
		xLength := mask.maskImage.Bounds().Dx()
		yLength := mask.maskImage.Bounds().Dy()
		if xIndex >= xLength || yIndex >= yLength {
			return false, errors.New("sorter: the mask lookup is out of the mask bounds")
		}

		color = utils.ColorToRgba(mask.maskImage.At(xIndex, yIndex))
	}

	// TODO: Refactor to use color.NRGBA
	_, _, l, _ := utils.RgbaToHsla(utils.ColorToRgba(color))

	if l == 0.0 {
		return true, nil
	} else if l == 1.0 {
		return false, nil
	} else {
		// TODO: The functioning of masks here is quite limited, we mask on a zero-one basis. During
		// interpolation, values creep in that do not correspond to this approach. Because of this,
		// we do a hack here, which is that we mask based on which value is closer to zero or one.
		if mask.isTranslated {
			if l < 0.5 {
				return true, nil
			} else {
				return false, nil
			}
		}

		return false, errors.New("sorter: the mask lookup found a invalid value and may be corrupted")
	}
}
