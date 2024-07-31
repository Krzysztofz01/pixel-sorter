package utils

// Perform a linear interpolation between value "a" and "b" and return the interpolated value at "t"
func Lerp(a, b, t float64) float64 {
	if t < 0.0 || t > 1.0 {
		panic("math-utils: the interpolated funtion argument must be between 0.0 and 1.0")
	}

	return (1.0-t)*a + t*b
}

// Clamp the x int between min and max. Return min if x is smaller than min and max if x is greater than max, otherwise return x.
func ClampInt(min, x, max int) int {
	if x < min {
		return min
	}

	if x > max {
		return max
	}

	return x
}

// Clamp the x float64 between min and max. Return min if x is smaller than min and max if x is greater than max, otherwise return x.
func ClampFloat64(min, x, max float64) float64 {
	if x < min {
		return min
	}

	if x > max {
		return max
	}

	return x
}

// Get the smallest float64 from the three provided values.
func Min3Float64(a, b, c float64) float64 {
	if a > b {
		a = b
	}

	if a > c {
		a = c
	}

	return a
}

// Get the greatest float64 from the three provided values.
func Max3Float64(a, b, c float64) float64 {
	if a < b {
		a = b
	}

	if a < c {
		a = c
	}

	return a
}

// Get the smallest uint8 from the two provided values.
func Min2Uint8(a, b uint8) uint8 {
	if a > b {
		return b
	}

	return a
}

// Get the greatest uint8 from the two provided values.
func Max2Uint8(a, b uint8) uint8 {
	if a < b {
		return b
	}

	return a
}
