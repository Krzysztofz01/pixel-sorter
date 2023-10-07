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
	cases := map[color.Color]bool{
		color.RGBA{0, 0, 0, 0}:          true,
		color.RGBA{0, 0, 0, 254}:        true,
		color.RGBA{0, 0, 0, 255}:        false,
		color.RGBA{255, 255, 255, 0}:    true,
		color.RGBA{255, 255, 255, 254}:  true,
		color.RGBA{255, 255, 255, 255}:  false,
		color.NRGBA{0, 0, 0, 0}:         true,
		color.NRGBA{0, 0, 0, 254}:       true,
		color.NRGBA{0, 0, 0, 255}:       false,
		color.NRGBA{255, 255, 255, 0}:   true,
		color.NRGBA{255, 255, 255, 254}: true,
		color.NRGBA{255, 255, 255, 255}: false,
		color.Black:                     false,
		color.White:                     false,
	}

	for color, expected := range cases {
		actual := HasAnyTransparency(color)

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

func TestShouldConvertColorToHslaComponents(t *testing.T) {
	cases := map[color.Color]struct {
		h int
		s float64
		l float64
		a float64
	}{
		color.RGBA{0, 0, 0, 0}:          {0, 0.0, 0.0, 0.0},
		color.RGBA{0, 0, 0, 255}:        {0, 0.0, 0.0, 1.0},
		color.RGBA{255, 255, 255, 0}:    {0, 0.0, 1.0, 0.0},
		color.RGBA{255, 255, 255, 255}:  {0, 0.0, 1.0, 1.0},
		color.RGBA{255, 0, 0, 0}:        {0, 1.0, 0.5, 0.0},
		color.RGBA{255, 0, 0, 255}:      {0, 1.0, 0.5, 1.0},
		color.RGBA{50, 100, 200, 0}:     {220, 0.60, 0.49, 0.0},
		color.RGBA{50, 100, 200, 255}:   {220, 0.60, 0.49, 1.0},
		color.RGBA{55, 99, 7, 255}:      {89, 0.87, 0.21, 1.0},
		color.NRGBA{0, 0, 0, 0}:         {0, 0.0, 0.0, 0.0},
		color.NRGBA{0, 0, 0, 255}:       {0, 0.0, 0.0, 1.0},
		color.NRGBA{255, 255, 255, 0}:   {0, 0.0, 0.0, 0.0},
		color.NRGBA{255, 255, 255, 255}: {0, 0.0, 1.0, 1.0},
		color.NRGBA{255, 0, 0, 0}:       {0, 0.0, 0.0, 0.0},
		color.NRGBA{255, 0, 0, 255}:     {0, 1.0, 0.5, 1.0},
		color.NRGBA{50, 100, 200, 0}:    {0, 0.0, 0.0, 0.0},
		color.NRGBA{50, 100, 200, 255}:  {220, 0.60, 0.49, 1.0},
		color.NRGBA{55, 99, 7, 255}:     {89, 0.87, 0.21, 1.0},
		color.Black:                     {0, 0.0, 0.0, 1.0},
		color.White:                     {0, 0.0, 1.0, 1.0},
	}

	const deltaHue float64 = 1.0
	const deltaSaturation float64 = 1e-2
	const deltaLightness float64 = 1e-2
	const deltaAlpha float64 = 1e-2

	for color, expected := range cases {
		hActual, sActual, lActual, aActual := ColorToHsla(color)

		if expected.s > 0.0 {
			assert.InDelta(t, expected.h, hActual, deltaHue)
		}

		assert.InDelta(t, expected.s, sActual, deltaSaturation)
		assert.InDelta(t, expected.l, lActual, deltaLightness)
		assert.InDelta(t, expected.a, aActual, deltaAlpha)
	}
}

func TestShouldBlendColorUsingLightenOnlyMode(t *testing.T) {
	cases := map[struct {
		a color.RGBA
		b color.RGBA
	}]color.RGBA{
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{0, 0, 0, 0xff}}:             {0, 0, 0, 0xff},
		{color.RGBA{255, 255, 255, 0xff}, color.RGBA{255, 255, 255, 0xff}}: {255, 255, 255, 0xff},
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{255, 255, 255, 0xff}}:       {255, 255, 255, 0xff},
		{color.RGBA{25, 50, 200, 0xff}, color.RGBA{200, 40, 20, 0xff}}:     {200, 50, 200, 0xff},
	}

	for c, exptected := range cases {
		actual := BlendRGBA(c.a, c.b, LightenOnly)

		assert.Equal(t, exptected, actual)
	}
}

func TestShouldBlendColorUsingDarkenOnlyMode(t *testing.T) {
	cases := map[struct {
		a color.RGBA
		b color.RGBA
	}]color.RGBA{
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{0, 0, 0, 0xff}}:             {0, 0, 0, 0xff},
		{color.RGBA{255, 255, 255, 0xff}, color.RGBA{255, 255, 255, 0xff}}: {255, 255, 255, 0xff},
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{255, 255, 255, 0xff}}:       {0, 0, 0, 0xff},
		{color.RGBA{25, 50, 200, 0xff}, color.RGBA{200, 40, 20, 0xff}}:     {25, 40, 20, 0xff},
	}

	for c, exptected := range cases {
		actual := BlendRGBA(c.a, c.b, DarkenOnly)

		assert.Equal(t, exptected, actual)
	}
}

func TestInterpolateColorShouldCorrectlyInterpolateColors(t *testing.T) {
	cases := map[struct {
		a color.Color
		b color.Color
		t float64
	}]color.Color{
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{255, 255, 255, 0xff}, 0.5}: color.RGBA{127, 127, 127, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.0}:   color.RGBA{255, 0, 0, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.25}:  color.RGBA{191, 0, 63, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.5}:   color.RGBA{127, 0, 127, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.75}:  color.RGBA{63, 0, 191, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 1.0}:   color.RGBA{0, 0, 255, 0xff},
	}

	for c, expected := range cases {
		actual := InterpolateColor(c.a, c.b, c.t)

		assert.Equal(t, expected, actual)
	}
}

func TestNrgbaToGrayscaleComponentShouldConvert(t *testing.T) {
	cases := map[color.NRGBA]int{
		{0, 0, 0, 0}:         0,
		{0, 0, 0, 255}:       0,
		{255, 255, 255, 0}:   255, // 0 if premultiplied
		{255, 255, 255, 255}: 255,
		{50, 100, 200, 0}:    96, // 0 if premultuplied
		{50, 100, 200, 255}:  96,
	}

	for color, expected := range cases {
		actual := NrgbaToGrayscaleComponent(color)

		assert.Equal(t, expected, actual)
	}
}

func TestRgbaToHslaShouldConvert(t *testing.T) {
	cases := map[color.RGBA]struct {
		h int
		s float64
		l float64
		a float64
	}{
		{0, 0, 0, 0}:         {0, 0.0, 0.0, 0.0},
		{0, 0, 0, 255}:       {0, 0.0, 0.0, 1.0},
		{255, 255, 255, 0}:   {0, 0.0, 1.0, 0.0},
		{255, 255, 255, 255}: {0, 0.0, 1.0, 1.0},
		{255, 0, 0, 0}:       {0, 1.0, 0.5, 0.0},
		{255, 0, 0, 255}:     {0, 1.0, 0.5, 1.0},
		{50, 100, 200, 0}:    {220, 0.60, 0.49, 0.0},
		{50, 100, 200, 255}:  {220, 0.60, 0.49, 1.0},
		{55, 99, 7, 255}:     {89, 0.87, 0.21, 1.0},
		{255, 250, 250, 0}:   {0, 1.0, 0.99, 0.0},
		{255, 250, 250, 255}: {0, 1.0, 0.99, 1.0},
		{255, 255, 0, 0}:     {60, 1.0, 0.5, 0.0},
		{255, 255, 0, 255}:   {60, 1.0, 0.5, 1.0},
	}

	const deltaHue float64 = 1.0
	const deltaSaturation float64 = 1e-2
	const deltaLightness float64 = 1e-2
	const deltaAlpha float64 = 1e-2

	for color, expected := range cases {
		hActual, sActual, lActual, aActual := RgbaToHsla(color)

		if expected.s > 0.0 {
			assert.InDelta(t, expected.h, hActual, deltaHue)
		}

		assert.InDelta(t, expected.s, sActual, deltaSaturation)
		assert.InDelta(t, expected.l, lActual, deltaLightness)
		assert.InDelta(t, expected.a, aActual, deltaAlpha)
	}
}

func TestBlendNrgbaShouldBlendColorsUsingLightenOnlyMode(t *testing.T) {
	cases := map[struct {
		a color.NRGBA
		b color.NRGBA
	}]color.NRGBA{
		{color.NRGBA{0, 0, 0, 0xff}, color.NRGBA{0, 0, 0, 0xff}}:             {0, 0, 0, 0xff},
		{color.NRGBA{255, 255, 255, 0xff}, color.NRGBA{255, 255, 255, 0xff}}: {255, 255, 255, 0xff},
		{color.NRGBA{0, 0, 0, 0xff}, color.NRGBA{255, 255, 255, 0xff}}:       {255, 255, 255, 0xff},
		{color.NRGBA{25, 50, 200, 0xff}, color.NRGBA{200, 40, 20, 0xff}}:     {200, 50, 200, 0xff},
	}

	for c, exptected := range cases {
		actual := BlendNrgba(c.a, c.b, LightenOnly)

		assert.Equal(t, exptected, actual)
	}
}

func TestBlendNrgbaShouldPanicForInvalidBlendingMode(t *testing.T) {
	assert.Panics(t, func() {
		c := color.NRGBA{}
		BlendNrgba(c, c, -1)
	})
}

func TestBlendNrgbaShouldBlendColorsUsingDarkenOnlyMode(t *testing.T) {
	cases := map[struct {
		a color.NRGBA
		b color.NRGBA
	}]color.NRGBA{
		{color.NRGBA{0, 0, 0, 0xff}, color.NRGBA{0, 0, 0, 0xff}}:             {0, 0, 0, 0xff},
		{color.NRGBA{255, 255, 255, 0xff}, color.NRGBA{255, 255, 255, 0xff}}: {255, 255, 255, 0xff},
		{color.NRGBA{0, 0, 0, 0xff}, color.NRGBA{255, 255, 255, 0xff}}:       {0, 0, 0, 0xff},
		{color.NRGBA{25, 50, 200, 0xff}, color.NRGBA{200, 40, 20, 0xff}}:     {25, 40, 20, 0xff},
	}

	for c, exptected := range cases {
		actual := BlendNrgba(c.a, c.b, DarkenOnly)

		assert.Equal(t, exptected, actual)
	}
}

func TestInterpolateRgbaShouldCorrectlyInterpolateColors(t *testing.T) {
	cases := map[struct {
		a color.RGBA
		b color.RGBA
		t float64
	}]color.Color{
		{color.RGBA{0, 0, 0, 0xff}, color.RGBA{255, 255, 255, 0xff}, 0.5}: color.RGBA{127, 127, 127, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.0}:   color.RGBA{255, 0, 0, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.25}:  color.RGBA{191, 0, 63, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.5}:   color.RGBA{127, 0, 127, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 0.75}:  color.RGBA{63, 0, 191, 0xff},
		{color.RGBA{255, 0, 0, 0xff}, color.RGBA{0, 0, 255, 0xff}, 1.0}:   color.RGBA{0, 0, 255, 0xff},
	}

	for c, expected := range cases {
		actual := InterpolateRgba(c.a, c.b, c.t)

		assert.Equal(t, expected, actual)
	}
}
