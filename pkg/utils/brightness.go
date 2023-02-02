package utils

import (
	"image/color"
	"math"
)

// Calculate the perceived brightness of the color given in the RGBA color space.
// https://stackoverflow.com/questions/596216/formula-to-determine-perceived-brightness-of-rgb-color
func CalculatePerceivedBrightness(c color.RGBA) float64 {
	luminance := calculateLuminance(c)

	if luminance <= 216.0/24389.0 {
		return (luminance * (24389.0 / 27.0)) / 100.0
	} else {
		return (math.Pow(luminance, 1.0/3.0)*116.0 - 16.0) / 100.0
	}
}

// Calculate the luminance of a color given in the RGBA color space. Helper function used to calculate the perceived brightness.
func calculateLuminance(c color.RGBA) float64 {
	rNorm, gNorm, bNorm := RgbaToNormalizedComponents(c)

	rLinear := calculateRgbComponentLinearValue(rNorm)
	gLinear := calculateRgbComponentLinearValue(gNorm)
	bLinear := calculateRgbComponentLinearValue(bNorm)

	return rLinear*0.2126 + gLinear*0.7152 + bLinear*0.0722
}

// Calculate the linear value of a given RGB component. Helper function used to calculate the RGB luminance
func calculateRgbComponentLinearValue(component float64) float64 {
	if component <= 0.04045 {
		return (component / 12.92)
	} else {
		return math.Pow((component+0.055)/1.055, 2.4)
	}
}
