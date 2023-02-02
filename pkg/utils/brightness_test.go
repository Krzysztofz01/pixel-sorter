package utils

import (
	"image/color"
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
