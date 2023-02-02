package utils

import (
	"errors"
	"image/color"
)

// Convert the color.RGBA struct to individual RGB components represented as integers in range from 0 to 255
func RgbaToIntComponents(c color.RGBA) (int, int, int) {
	r32, g32, b32, _ := c.RGBA()
	r := int(r32 >> 8)
	g := int(g32 >> 8)
	b := int(b32 >> 8)

	return r, g, b
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
