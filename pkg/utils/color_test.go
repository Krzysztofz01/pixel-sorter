package utils

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldConvertColorToGrayscaleComponent(t *testing.T) {
	cases := map[color.Color]int{
		color.RGBA{0, 0, 0, 0}:          0,
		color.RGBA{0, 0, 0, 255}:        0,
		color.RGBA{255, 255, 255, 0}:    255,
		color.RGBA{255, 255, 255, 255}:  255,
		color.RGBA{50, 100, 200, 0}:     96,
		color.RGBA{50, 100, 200, 255}:   96,
		color.NRGBA{0, 0, 0, 0}:         0,
		color.NRGBA{0, 0, 0, 255}:       0,
		color.NRGBA{255, 255, 255, 0}:   0,
		color.NRGBA{255, 255, 255, 255}: 255,
		color.NRGBA{50, 100, 200, 0}:    0,
		color.NRGBA{50, 100, 200, 255}:  96,
		color.White:                     255,
		color.Black:                     0,
	}

	for color, expected := range cases {
		actual := ColorToGrayscaleComponent(color)

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

func TestShouldConvertColorToRgba(t *testing.T) {
	cases := map[color.Color]color.RGBA{
		color.RGBA{0, 0, 0, 0}:          {0, 0, 0, 0},
		color.RGBA{0, 0, 0, 254}:        {0, 0, 0, 254},
		color.RGBA{0, 0, 0, 255}:        {0, 0, 0, 255},
		color.RGBA{255, 255, 255, 0}:    {255, 255, 255, 0},
		color.RGBA{255, 255, 255, 254}:  {255, 255, 255, 254},
		color.RGBA{255, 255, 255, 255}:  {255, 255, 255, 255},
		color.NRGBA{0, 0, 0, 0}:         {0, 0, 0, 0},
		color.NRGBA{0, 0, 0, 254}:       {0, 0, 0, 254},
		color.NRGBA{0, 0, 0, 255}:       {0, 0, 0, 255},
		color.NRGBA{255, 255, 255, 0}:   {0, 0, 0, 0},
		color.NRGBA{255, 255, 255, 254}: {254, 254, 254, 254},
		color.NRGBA{255, 255, 255, 255}: {255, 255, 255, 255},
		color.Gray{Y: 244}:              {244, 244, 244, 255},
		color.White:                     {255, 255, 255, 255},
		color.Black:                     {0, 0, 0, 255},
	}

	for color, expected := range cases {
		actual := ColorToRgba(color)

		assert.Equal(t, expected, actual)
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

func TestShouldBlendColorUsingLightenOnlyMode(t *testing.T) {
	// TODO: Implement more test cases
	aColor := color.RGBA{25, 50, 200, 0xff}
	bColor := color.RGBA{200, 40, 20, 0xff}

	expected := color.RGBA{200, 50, 200, 0xff}
	actual := BlendRGBA(aColor, bColor, LightenOnly)

	assert.Equal(t, expected, actual)
}

func TestShouldBlendColorUsingDarkenOnlyMode(t *testing.T) {
	// TODO: Implement more test cases
	aColor := color.RGBA{25, 50, 200, 0xff}
	bColor := color.RGBA{200, 40, 20, 0xff}

	expected := color.RGBA{25, 40, 20, 0xff}
	actual := BlendRGBA(aColor, bColor, DarkenOnly)

	assert.Equal(t, expected, actual)
}
