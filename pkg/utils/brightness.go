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
		return (luminanceRangeCubeRoot(luminance)*116.0 - 16.0) / 100.0
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
		// TODO: Implement custom function for this, just like with luminance cbrt
		return math.Pow((component+0.055)/1.055, 2.4)
	}
}

// A third attempt to optimize the linear RGB component calculation optimization via approximation. A polynomial function with 5 roots has been
// created using polynomial regresion.The value from the calculated via the polynomial is further approximated using "Newton-Raphson method" (3 iterations)
// This method, despite the efforts, is still not precise enough and has been left here mainly as a curiosity.
func rgbComponentLinearApproximation(component float64) float64 {
	if component <= 0.04045 {
		return (component / 12.92)
	}

	polynomial := func(x float64) float64 {
		x2 := x * x
		return (0.054078451133248 * x2 * x2 * x) - (0.224682405485073 * x2 * x2) + (0.594031158357391 * x2 * x) + (0.546897380766604 * x2) + (0.028604496057605 * x) + 0.001100839502532
	}

	derivative := func(x float64) float64 {
		x2 := x * x
		return (0.27039225566624 * x2 * x2) - (0.898729621940292 * x2 * x) + (1.78209347507217 * x2) + (1.09379476153321 * x) + 0.028604496057605
	}

	value := polynomial(component)
	for i := 0; i < 3; i += 1 {
		value = value - (polynomial(value) / derivative(value))
	}

	return value
}

// The profiling showed that the brightness calculation is very slow beacuse of the cubic root operation on the luminance value.
// The analysis of the brightness calculation showed, that the cube root results are always in the range: 4/29 <= cbrt(luminance) <= 1
// We are using the "Newton-Raphson method" (3 iterations) to perform the cube root approximation. In order to get a precise initial value we are
// using a square polynomial (there is also a linear one...) with a formula obtained from a polynomial regression performed on the
// cube root X and Y values, where the X is in range between 4/29 and 1.
func luminanceRangeCubeRoot(x float64) float64 {
	// Cube root square polynomial approximation formula
	reg := (-0.358955950652834 * x * x) + (0.934309346877746 * x) + 0.414814427166639

	// Cube root linear approximation formula
	// req := 0.525842230617626*x + 0.508563748107305

	for i := 0; i < 3; i += 1 {
		regp2 := reg * reg
		reg = reg - ((regp2*reg)-x)/(3*regp2)
	}

	return reg
}
