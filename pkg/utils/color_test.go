package utils

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldConvertRGBAToRgbComponents(t *testing.T) {
	cases := map[color.RGBA]struct{ r, g, b int }{
		{0, 0, 0, 0}:         {0, 0, 0},
		{0, 0, 0, 255}:       {0, 0, 0},
		{255, 255, 255, 0}:   {255, 255, 255},
		{255, 255, 255, 255}: {255, 255, 255},
		{50, 100, 200, 0}:    {50, 100, 200},
		{50, 100, 200, 255}:  {50, 100, 200},
	}

	for rgba, expected := range cases {
		rActual, gActual, bActual := RgbaToIntComponents(rgba)

		assert.Equal(t, expected.r, rActual)
		assert.Equal(t, expected.g, gActual)
		assert.Equal(t, expected.b, bActual)
	}
}

func TestShouldConvertNRGBAToRgbComponents(t *testing.T) {
	cases := map[color.NRGBA]struct{ r, g, b int }{
		{0, 0, 0, 0}:         {0, 0, 0},
		{0, 0, 0, 255}:       {0, 0, 0},
		{255, 255, 255, 0}:   {255, 255, 255},
		{255, 255, 255, 255}: {255, 255, 255},
		{50, 100, 200, 0}:    {50, 100, 200},
		{50, 100, 200, 255}:  {50, 100, 200},
	}

	for nrgba, expected := range cases {
		rActual, gActual, bActual := NrgbaToIntComponents(nrgba)

		assert.Equal(t, expected.r, rActual)
		assert.Equal(t, expected.g, gActual)
		assert.Equal(t, expected.b, bActual)
	}
}

func TestShouldConvertRGBAToGrayscaleComponent(t *testing.T) {
	cases := map[color.RGBA]int{
		{0, 0, 0, 0}:         0,
		{0, 0, 0, 255}:       0,
		{255, 255, 255, 0}:   255,
		{255, 255, 255, 255}: 255,
		{50, 100, 200, 0}:    96,
		{50, 100, 200, 255}:  96,
	}

	for rgba, expected := range cases {
		actual := RgbaToGrayscaleComponent(rgba)

		assert.Equal(t, expected, actual)
	}
}

func TestShouldConvertNRGBAToGrayscaleComponent(t *testing.T) {
	cases := map[color.NRGBA]int{
		{0, 0, 0, 0}:         0,
		{0, 0, 0, 255}:       0,
		{255, 255, 255, 0}:   255,
		{255, 255, 255, 255}: 255,
		{50, 100, 200, 0}:    96,
		{50, 100, 200, 255}:  96,
	}

	for nrgba, expected := range cases {
		actual := NrgbaToGrayscaleComponent(nrgba)

		assert.Equal(t, expected, actual)
	}
}

func TestShouldConvertRGBAToNormalizedRgbComponents(t *testing.T) {
	cases := map[color.RGBA]struct{ r, g, b float64 }{
		{0, 0, 0, 0}:         {0.0, 0.0, 0.0},
		{0, 0, 0, 255}:       {0.0, 0.0, 0.0},
		{255, 255, 255, 0}:   {1.0, 1.0, 1.0},
		{255, 255, 255, 255}: {1.0, 1.0, 1.0},
	}

	const delta float64 = 1e-7

	for rgba, expected := range cases {
		rActual, gActual, bActual := RgbaToNormalizedComponents(rgba)

		assert.InDelta(t, expected.r, rActual, delta)
		assert.InDelta(t, expected.g, gActual, delta)
		assert.InDelta(t, expected.b, bActual, delta)
	}
}

func TestShouldTellIfRGBAHasAnyTransparency(t *testing.T) {
	cases := map[color.RGBA]bool{
		{0, 0, 0, 0}:         true,
		{0, 0, 0, 254}:       true,
		{0, 0, 0, 255}:       false,
		{255, 255, 255, 0}:   true,
		{255, 255, 255, 254}: true,
		{255, 255, 255, 255}: false,
	}

	for rgba, expected := range cases {
		actual := HasAnyTransparency(rgba)

		assert.Equal(t, expected, actual)
	}
}

func TestShouldConvertColorToRgbaIfColorIsRgba(t *testing.T) {
	cases := map[color.Color]color.RGBA{
		color.RGBA{0, 0, 0, 0}:         {0, 0, 0, 0},
		color.RGBA{0, 0, 0, 254}:       {0, 0, 0, 254},
		color.RGBA{0, 0, 0, 255}:       {0, 0, 0, 255},
		color.RGBA{255, 255, 255, 0}:   {255, 255, 255, 0},
		color.RGBA{255, 255, 255, 254}: {255, 255, 255, 254},
		color.RGBA{255, 255, 255, 255}: {255, 255, 255, 255},
	}

	for color, expected := range cases {
		actual, err := ColorToRgba(color)

		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	}
}

func TestShouldNotConvertColorToRgbaIfColorIsNotRgba(t *testing.T) {
	cases := map[color.Color]color.RGBA{
		color.NRGBA{0, 0, 0, 0}:         {0, 0, 0, 0},
		color.NRGBA{0, 0, 0, 254}:       {0, 0, 0, 254},
		color.NRGBA{0, 0, 0, 255}:       {0, 0, 0, 255},
		color.NRGBA{255, 255, 255, 0}:   {255, 255, 255, 0},
		color.NRGBA{255, 255, 255, 254}: {255, 255, 255, 254},
		color.NRGBA{255, 255, 255, 255}: {255, 255, 255, 255},
	}

	for color := range cases {
		_, err := ColorToRgba(color)
		assert.Error(t, err)
	}
}

func TestShouldConvertRgbaToHslComponents(t *testing.T) {
	cases := map[color.RGBA]struct {
		h int
		s float64
		l float64
	}{
		{0, 0, 0, 0}:         {0, 0.0, 0.0},
		{0, 0, 0, 255}:       {0.0, 0.0, 0.0},
		{255, 255, 255, 0}:   {0.0, 0.0, 1.0},
		{255, 255, 255, 255}: {0.0, 0.0, 1.0},
		{50, 100, 200, 0}:    {220, 0.60, 0.49},
		{50, 100, 200, 255}:  {220, 0.60, 0.49},
	}

	const deltaHue float64 = 1.0
	const deltaSaturation float64 = 0.1
	const deltaLightness float64 = 0.1

	for rgba, expected := range cases {
		hActual, sActual, lActual := RgbaToHsl(rgba)

		if expected.s > 0.0 {
			assert.InDelta(t, expected.h, hActual+1, deltaHue)
		}

		assert.InDelta(t, expected.s, sActual, deltaSaturation)
		assert.InDelta(t, expected.l, lActual, deltaLightness)
	}
}
