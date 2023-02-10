package sorter

import (
	"errors"
	"fmt"
	"image"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

// A structure representing a wrapper over a image that works as a mask
type Mask struct {
	maskImage image.Image
}

// TODO: Place the image/mask size validation here?

// Create a new mask instance from a given image
func CreateMask(i image.Image) (*Mask, error) {
	// TODO: It is not efficient to redraw the whole mask, and we dont need it in a drawable format, but we
	// are doing it to ensure the mask is a image.RGBA
	drawableMask, err := utils.GetDrawableImage(i)
	if err != nil {
		return nil, fmt.Errorf("sorter: failed the convert the mask to drawable version: %w", err)
	}

	xLength := drawableMask.Bounds().Dx()
	yLength := drawableMask.Bounds().Dy()

	for xIndex := 0; xIndex < xLength; xIndex += 1 {
		for yIndex := 0; yIndex < yLength; yIndex += 1 {
			color, err := utils.ColorToRgba(drawableMask.At(xIndex, yIndex))
			if err != nil {
				return nil, fmt.Errorf("sorter: failed to convert the color in the mask validation process: %w", err)
			}

			_, s, l := utils.RgbaToHsl(color)

			if s != 0.0 || (l != 0.0 && l != 1.0) {
				return nil, errors.New("sorter: the mask contains a invalid color")
			}
		}
	}

	mask := new(Mask)
	mask.maskImage = drawableMask
	return mask, nil
}

// Perform a mask lookup to check if the mask is masking at the given location
func (mask *Mask) IsMasked(xIndex, yIndex int) (bool, error) {
	xLength := mask.maskImage.Bounds().Dx()
	yLength := mask.maskImage.Bounds().Dy()
	if xIndex >= xLength || yIndex >= yLength {
		return false, errors.New("sorter: the mask lookup is out of the mask bounds")
	}

	color, err := utils.ColorToRgba(mask.maskImage.At(xIndex, yIndex))
	if err != nil {
		return false, fmt.Errorf("sorter: failed to convert the color in the mask lookup process: %w", err)
	}

	_, _, l := utils.RgbaToHsl(color)

	if l == 0.0 {
		return true, nil
	} else if l == 1.0 {
		return false, nil
	} else {
		return false, errors.New("sorter: the mask lookup found a invalid value and may be corrupted")
	}
}
