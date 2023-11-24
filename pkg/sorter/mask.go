package sorter

import (
	"errors"
	"fmt"
	"image"

	"github.com/Krzysztofz01/pimit"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

type Mask interface {
	// Return a byte value representing the mask grayscale value at the given position represented by the x and y position.
	At(x, y int) (uint8, error)

	// Return a boolean value representing if the mask is masking at the given position represented by the x and y postition (value < 127 is masked).
	AtB(x, y int) (bool, error)

	// Return a byte value representing the mask grayscale value at the given position represented by the i index.
	AtByIndex(i int) (uint8, error)

	// Return a boolean value representing if the mask is masking at the given position represented by the i index (value < 127 is masked).
	AtByIndexB(i int) (bool, error)
}

type mask struct {
	maskBuffer []uint8
	width      int
	isEmpty    bool
}

func CreateMaskFromNrgba(i *image.NRGBA) (Mask, error) {
	if i == nil {
		return nil, errors.New("mask: the provided mask image reference is nil")
	}

	mask := &mask{
		maskBuffer: make([]uint8, i.Bounds().Dx()*i.Bounds().Dy()),
		width:      i.Bounds().Dx(),
		isEmpty:    false,
	}

	errt := utils.NewErrorTrap()
	pimit.ParallelNrgbaRead(i, func(x, y int, r, g, b, a uint8) {
		if r != g || r != b || g != b {
			errt.Set(fmt.Errorf("mask: the mask image contains a invalid color at x=%d y=%d", x, y))
			return
		}

		index := y*i.Bounds().Dx() + x
		mask.maskBuffer[index] = r
	})

	if errt.IsSet() {
		return nil, errt.Err()
	}

	return mask, nil
}

func CreateMaskFromRgba(i *image.RGBA) (Mask, error) {
	if i == nil {
		return nil, errors.New("mask: the provided mask image reference is nil")
	}

	mask := &mask{
		maskBuffer: make([]uint8, i.Bounds().Dx()*i.Bounds().Dy()),
		width:      i.Bounds().Dx(),
		isEmpty:    false,
	}

	errt := utils.NewErrorTrap()
	pimit.ParallelRgbaRead(i, func(x, y int, r, g, b, a uint8) {
		if r != g || r != b || g != b {
			errt.Set(fmt.Errorf("mask: the mask image contains a invalid color at x=%d y=%d", x, y))
			return
		}

		index := y*i.Bounds().Dx() + x
		mask.maskBuffer[index] = r
	})

	if errt.IsSet() {
		return nil, errt.Err()
	}

	return mask, nil
}

func CreateEmptyMask() Mask {
	return &mask{
		maskBuffer: nil,
		width:      0,
		isEmpty:    true,
	}
}

func (m *mask) At(x, y int) (uint8, error) {
	return m.AtByIndex(y*m.width + x)
}

func (m *mask) AtB(x, y int) (bool, error) {
	return m.AtByIndexB(y*m.width + x)
}

func (m *mask) AtByIndex(i int) (uint8, error) {
	if m.isEmpty {
		return 0xff, nil
	}

	if i < 0 || i >= len(m.maskBuffer) {
		return 0x00, errors.New("mask: the specified index is out of the mask range")
	}

	return m.maskBuffer[i], nil
}

func (m *mask) AtByIndexB(i int) (bool, error) {
	at, err := m.AtByIndex(i)
	if err != nil {
		return false, err
	}

	if at < 127 {
		return true, nil
	} else {
		return false, nil
	}
}
