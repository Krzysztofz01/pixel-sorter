package sorter

import (
	"errors"
	"fmt"
	"image"
	"image/color"

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
	return CreateImageMask(mImg, mImg.Bounds(), 0)
}

// Create a new mask instance from a given image and bounds of the image to be masked. The trnslateAngle parameter
// indicates whether the lookup mask should be interpreted from a given angle.
func CreateImageMask(mImg image.Image, targetImageBounds image.Rectangle, translateAngle int) (*Mask, error) {
	if mImg.Bounds().Dx() != targetImageBounds.Dx() || mImg.Bounds().Dy() != targetImageBounds.Dy() {
		return nil, errors.New("sorter: mask image and target image sizes are not matching")
	}

	// TODO: It is not efficient to redraw the whole mask, and we dont need it in a drawable format, but we
	// are doing it to ensure the mask is a image.RGBA
	drawableMask, err := utils.GetDrawableImage(mImg)
	if err != nil {
		return nil, fmt.Errorf("sorter: failed the convert the mask to drawable version: %w", err)
	}

	xLength := drawableMask.Bounds().Dx()
	yLength := drawableMask.Bounds().Dy()

	for xIndex := 0; xIndex < xLength; xIndex += 1 {
		for yIndex := 0; yIndex < yLength; yIndex += 1 {
			color := utils.ColorToRgba(drawableMask.At(xIndex, yIndex))
			_, s, l := utils.RgbaToHsl(color)

			if s != 0.0 || (l != 0.0 && l != 1.0) {
				return nil, errors.New("sorter: the mask contains a invalid color")
			}
		}
	}

	mask := new(Mask)
	mask.maskImage = drawableMask
	mask.isEmpty = false

	if translateAngle != 0 {
		mask.maskImageTranslated = utils.RotateImage(drawableMask, translateAngle)
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

	_, _, l := utils.RgbaToHsl(utils.ColorToRgba(color))

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
