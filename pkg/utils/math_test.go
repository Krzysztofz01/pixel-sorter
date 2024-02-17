package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLerpShouldPanicForInvalidInterpolatedFunctionArgument(t *testing.T) {
	cases := []float64{
		-1.0,
		-0.1,
		1.1,
		2.0,
	}

	for _, tValue := range cases {
		assert.Panics(t, func() {
			Lerp(0.0, 1.0, tValue)
		})
	}
}

func TestLerpShouldCorrecltyCalculateTheLinearInterpolation(t *testing.T) {
	cases := map[struct {
		a float64
		b float64
		t float64
	}]float64{
		{0.0, 1.0, 0.5}: 0.5,
		{0.0, 1.0, 0.0}: 0.0,
		{0.0, 1.0, 1.0}: 1.0,
		{2, 102, 0.5}:   52.0,
		{2, 102, 0.0}:   2.0,
		{2, 102, 1.0}:   102.0,
	}

	var delta float64 = 1e-7

	for c, expected := range cases {
		actual := Lerp(c.a, c.b, c.t)

		assert.InDelta(t, expected, actual, delta)
	}
}
