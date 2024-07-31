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

func TestClampIntShouldCorrectlyClampValues(t *testing.T) {
	cases := map[struct {
		min int
		x   int
		max int
	}]int{
		{2, 0, 6}: 2,
		{2, 2, 6}: 2,
		{2, 4, 6}: 4,
		{2, 6, 6}: 6,
		{2, 8, 6}: 6,
	}

	for c, expected := range cases {
		actual := ClampInt(c.min, c.x, c.max)

		assert.Equal(t, expected, actual)
	}
}

func TestClampFloat64ShouldCorrectlyClampValues(t *testing.T) {
	cases := map[struct {
		min float64
		x   float64
		max float64
	}]float64{
		{2, 0, 6}: 2,
		{2, 2, 6}: 2,
		{2, 4, 6}: 4,
		{2, 6, 6}: 6,
		{2, 8, 6}: 6,
	}

	var delta float64 = 1e-7

	for c, expected := range cases {
		actual := ClampFloat64(c.min, c.x, c.max)

		assert.InDelta(t, expected, actual, delta)
	}
}

func TestMin3Float64ShouldCorrecltyReturnMinValue(t *testing.T) {
	cases := map[struct {
		a float64
		b float64
		c float64
	}]float64{
		{1, 2, 3}: 1,
		{3, 1, 2}: 1,
		{2, 3, 1}: 1,
		{1, 1, 1}: 1,
	}

	var delta float64 = 1e-7

	for c, expected := range cases {
		actual := Min3Float64(c.a, c.b, c.c)

		assert.InDelta(t, expected, actual, delta)
	}
}

func TestMax3Float64ShouldCorrecltyReturnMaxValue(t *testing.T) {
	cases := map[struct {
		a float64
		b float64
		c float64
	}]float64{
		{1, 2, 3}: 3,
		{3, 1, 2}: 3,
		{2, 3, 1}: 3,
		{1, 1, 1}: 1,
	}

	var delta float64 = 1e-7

	for c, expected := range cases {
		actual := Max3Float64(c.a, c.b, c.c)

		assert.InDelta(t, expected, actual, delta)
	}
}

func TestMin2Uint8ShouldCorrecltyReturnMinValue(t *testing.T) {
	cases := map[struct {
		a uint8
		b uint8
	}]uint8{
		{1, 2}: 1,
		{3, 1}: 1,
		{1, 1}: 1,
	}

	var delta float64 = 1e-7

	for c, expected := range cases {
		actual := Min2Uint8(c.a, c.b)

		assert.InDelta(t, expected, actual, delta)
	}
}

func TestMax2Uint8ShouldCorrecltyReturnMaxValue(t *testing.T) {
	cases := map[struct {
		a uint8
		b uint8
	}]uint8{
		{1, 2}: 2,
		{3, 1}: 3,
		{1, 1}: 1,
	}

	var delta float64 = 1e-7

	for c, expected := range cases {
		actual := Max2Uint8(c.a, c.b)

		assert.InDelta(t, expected, actual, delta)
	}
}
