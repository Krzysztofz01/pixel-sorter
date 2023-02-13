package utils

import (
	"image/color"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldCalculateCorrectPerceivedBrightness(t *testing.T) {
	cases := map[color.RGBA]float64{
		{0, 0, 0, 0}:         0.0,
		{0, 0, 0, 255}:       0.0,
		{255, 255, 255, 0}:   1.0,
		{255, 255, 255, 255}: 1.0,
	}

	const delta float64 = 1e-7

	for rgba, expected := range cases {
		actual := CalculatePerceivedBrightness(rgba)

		assert.InDelta(t, expected, actual, delta)
	}
}

func TestLuminanceRangeCubicRootShouldCalculatePreciseValuesInCorrectRange(t *testing.T) {
	const (
		min        = 16.0 / 116.0
		max        = 1.0
		iterations = 200
		step       = (max - min) / float64(iterations)
		delta      = 1e-7
	)

	for x := min; x < max; x += step {
		expected := math.Cbrt(x)
		actual := luminanceRangeCubeRoot(x)

		assert.InDelta(t, expected, actual, delta)
	}
}
