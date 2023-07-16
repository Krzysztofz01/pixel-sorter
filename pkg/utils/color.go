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

// Convert the color.Color interface instance to the Y grayscale component represented as a integer in range from 0 to 255
func ColorToGrayscaleComponent(c color.Color) int {
	rgba := ColorToRgba(c)

	y := (float64(rgba.R) * 0.299) + (float64(rgba.G) * 0.587) + (float64(rgba.B) * 0.114)
	return int(math.Max(0, math.Min(255, y)))
}

// Convert the color.RGBA struct tu individual RGB components represented as floating point numbers in range from 0.0 to 1.0
func RgbaToNormalizedComponents(c color.RGBA) (float64, float64, float64) {
	rNorm := float64(c.R) / 255.0
	gNorm := float64(c.G) / 255.0
	bNorm := float64(c.B) / 255.0

	return rNorm, gNorm, bNorm
}

// Return a boolean value indicating if the given color.Color interface implementation has the alpha channel <255
func HasAnyTransparency(c color.Color) bool {
	_, _, _, a32 := c.RGBA()

	return int(a32>>8) < 255
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

// Convert a color represented as color.Color interface implementation to HSL+Alpha components where Hue is expressed in degress (0-360) and the
// saturation, lightnes and alpha in percentage (0.0-1.0)
func ColorToHsla(c color.Color) (int, float64, float64, float64) {
	rgba := ColorToRgba(c)
	rNorm, gNorm, bNorm := RgbaToNormalizedComponents(rgba)

	alpha := float64(rgba.A) / 255.0

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

	return hue, saturation, lightness, alpha
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

// Create a average color based on the provided colors depending
// on the RGBA color components.
// TODO: Unit tests implementation
func AverageColor(c ...color.Color) color.Color {
	r, g, b, a := 0, 0, 0, 0
	for _, color := range c {
		rgba := ColorToRgba(color)

		r += int(rgba.R)
		g += int(rgba.G)
		b += int(rgba.B)
		a += int(rgba.A)
	}

	paramCount := len(c)
	return color.RGBA{
		R: uint8(r / paramCount),
		G: uint8(g / paramCount),
		B: uint8(b / paramCount),
		A: uint8(a / paramCount),
	}
}
