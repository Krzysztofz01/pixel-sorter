package utils

import (
	"errors"
	"image/color"
	"math"
)

// Convert the color.RGBA struct to individual RGB components represented as integers in range from 0 to 255
func RgbaToIntComponents(c color.RGBA) (int, int, int) {
	r32, g32, b32, _ := c.RGBA()
	r := int(r32 >> 8)
	g := int(g32 >> 8)
	b := int(b32 >> 8)

	return r, g, b
}

// Convert the color.NRGBA struct to individual RGB components represented as integers in range from 0 to 255
func NrgbaToIntComponents(c color.NRGBA) (int, int, int) {
	r := int(c.R)
	g := int(c.G)
	b := int(c.B)

	return r, g, b
}

// Convert the color.RGBA struct to the Y grayscale component represented as integer in range from 0 to 255
func RgbaToGrayscaleComponent(c color.RGBA) int {
	r, g, b := RgbaToIntComponents(c)

	y := (float64(r) * 0.299) + (float64(g) * 0.587) + (float64(b) * 0.114)
	return int(math.Min(255, math.Max(0, y)))
}

// Convert the color.NRGBA struct to the Y grayscale component represented as integer in range from 0 to 255
func NrgbaToGrayscaleComponent(c color.NRGBA) int {
	r, g, b := NrgbaToIntComponents(c)

	y := (float64(r) * 0.299) + (float64(g) * 0.587) + (float64(b) * 0.114)
	return int(math.Min(255, math.Max(0, y)))
}

// Convert the color.RGBA struct tu individual RGB components represented as floating point numbers in range from 0.0 to 1.0
func RgbaToNormalizedComponents(c color.RGBA) (float64, float64, float64) {
	r, g, b := RgbaToIntComponents(c)
	rNorm := float64(r) / 255.0
	gNorm := float64(g) / 255.0
	bNorm := float64(b) / 255.0

	return rNorm, gNorm, bNorm
}

// Return a boolean value indicating if the given color.RGBA color has the alpha channel >255
func HasAnyTransparency(c color.RGBA) bool {
	_, _, _, a32 := c.RGBA()
	a := int(a32 >> 8)

	return a < 255
}

// Convert a color represented as color.Color interface to color.RGBA struct. This function will return an error if the underlying color is not a color.RGBA
func ColorToRgba(c color.Color) (color.RGBA, error) {
	rgba, ok := c.(color.RGBA)
	if !ok {
		return color.RGBA{}, errors.New("color-utils: conversion failed becuse the underlying color implementation is not RGBA")
	}

	return rgba, nil
}

// Convert a color represented as color.RGBA to HSL components where Hue is expressed in degress (0-360) and the saturation and lightnes in percentage (0.0-1.0)
func RgbaToHsl(c color.RGBA) (int, float64, float64) {
	rNorm, gNorm, bNorm := RgbaToNormalizedComponents(c)

	min := math.Min(rNorm, math.Min(gNorm, bNorm))
	max := math.Max(rNorm, math.Max(gNorm, bNorm))
	delta := max - min

	lightness := (max + min) / 2.0
	saturation := 0.0
	hue := 0

	if delta != 0.0 {
		if lightness <= 0.5 {
			saturation = delta / (max + min)
		} else {
			saturation = delta / (2.0 - max - min)
		}

		hueNorm := 0.0
		if max == rNorm {
			hueNorm = ((gNorm - bNorm) / 6.0) / delta
		} else if max == gNorm {
			hueNorm = (1.0 / 3.0) + ((bNorm-rNorm)/6.0)/delta
		} else {
			hueNorm = (2.0 / 3.0) + ((rNorm-gNorm)/6.0)/delta
		}

		if hueNorm < 0.0 {
			hueNorm += 1.0
		}

		if hueNorm > 1.0 {
			hueNorm -= 1.0
		}

		hue = int(math.Round(hueNorm * 360))
	}

	return hue, saturation, lightness
}
