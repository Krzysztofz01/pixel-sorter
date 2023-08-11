package utils

// Perform a linear interpolation between value "a" and "b" and return the interpolated value at "t"
func Lerp(a, b, t float64) float64 {
	if t < 0.0 || t > 1.0 {
		panic("math-utils: the interpolated funtion argument must be between 0.0 and 1.0")
	}

	return (1.0-t)*a + t*b
}
