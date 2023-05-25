package utils

import (
	"image/color"
	"math"
)

// Representation of a blending mode algorithm
type BlendingMode int

const (
	LightenOnly BlendingMode = iota
	DarkenOnly
)

// Convert the color.NRGBA struct to individual RGB components represented as integers in range from 0 to 255
func NrgbaToIntComponents(c color.NRGBA) (int, int, int) {
	r := int(c.R)
	g := int(c.G)
	b := int(c.B)

	return r, g, b
}

// Convert the color.RGBA struct to the Y grayscale component represented as integer in range from 0 to 255
func RgbaToGrayscaleComponent(c color.RGBA) int {
	y := (float64(c.R) * 0.299) + (float64(c.G) * 0.587) + (float64(c.B) * 0.114)
	return int(math.Min(255, math.Max(0, y)))
}

// Convert the color.NRGBA struct to the Y grayscale component represented as integer in range from 0 to 255
func NrgbaToGrayscaleComponent(c color.NRGBA) int {
	r, g, b := NrgbaToIntComponents(c)

	y := (float64(r) * 0.299) + (float64(g) * 0.587) + (float64(b) * 0.114)
	return int(math.Min(255, math.Max(0, y)))
}

// Convert the color.RGBA struct tu individual RGB components represented as floating point numbers in range from 0.0 to 1.0
func RgbaToNormalizedComponents(c color.RGBA) (float64, float64, float64) {
	rNorm := float64(c.R) / 255.0
	gNorm := float64(c.G) / 255.0
	bNorm := float64(c.B) / 255.0

	return rNorm, gNorm, bNorm
}

// Return a boolean value indicating if the given color.RGBA color has the alpha channel >255
func HasAnyTransparency(c color.RGBA) bool {
	_, _, _, a32 := c.RGBA()
	a := int(a32 >> 8)

	return a < 255
}

// Convert a color represented as color.Color interface to color.RGBA struct. If the underlying color is a color.RGBA the original struct
// will be returned, otherwise a new color.RGBA instance will be created
func ColorToRgba(c color.Color) color.RGBA {
	if rgba, ok := c.(color.RGBA); ok {
		return rgba
	}

	r32, g32, b32, a32 := c.RGBA()
	return color.RGBA{
		R: uint8(r32 >> 8),
		G: uint8(g32 >> 8),
		B: uint8(b32 >> 8),
		A: uint8(a32 >> 8),
	}
}

// Convert a color represented as color.RGBA to HSL components where Hue is expressed in degress (0-360) and the saturation and lightnes in percentage (0.0-1.0)
func RgbaToHsl(c color.RGBA) (int, float64, float64) {
	rNorm, gNorm, bNorm := RgbaToNormalizedComponents(c)

	min := math.Min(rNorm, math.Min(gNorm, bNorm))
	max := math.Max(rNorm, math.Max(gNorm, bNorm))
	delta := max - min

	lightness := (max + min) / 2.0
	saturation := 0.0
	hue := 0

	if delta != 0.0 {
		if lightness <= 0.5 {
			saturation = delta / (max + min)
		} else {
			saturation = delta / (2.0 - max - min)
		}

		hueNorm := 0.0
		if max == rNorm {
			hueNorm = ((gNorm - bNorm) / 6.0) / delta
		} else if max == gNorm {
			hueNorm = (1.0 / 3.0) + ((bNorm-rNorm)/6.0)/delta
		} else {
			hueNorm = (2.0 / 3.0) + ((rNorm-gNorm)/6.0)/delta
		}

		if hueNorm < 0.0 {
			hueNorm += 1.0
		}

		if hueNorm > 1.0 {
			hueNorm -= 1.0
		}

		hue = int(math.Round(hueNorm * 360))
	}

	return hue, saturation, lightness
}

// Perform blending of two colors according to a given blending mode
// TODO: Currently the alpha channel of the output color has a 0xff fixed value
func BlendRGBA(a, b color.RGBA, mode BlendingMode) color.RGBA {
	switch mode {
	case LightenOnly:
		{
			r := uint8(math.Max(float64(a.R), float64(b.R)))
			g := uint8(math.Max(float64(a.G), float64(b.G)))
			b := uint8(math.Max(float64(a.B), float64(b.B)))

			return color.RGBA{r, g, b, 0xff}
		}
	case DarkenOnly:
		{
			r := uint8(math.Min(float64(a.R), float64(b.R)))
			g := uint8(math.Min(float64(a.G), float64(b.G)))
			b := uint8(math.Min(float64(a.B), float64(b.B)))

			return color.RGBA{r, g, b, 0xff}
		}
	default:
		panic("color-utils: undefined blending mode provided")
	}
}
