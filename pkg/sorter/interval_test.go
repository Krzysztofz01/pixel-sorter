package sorter

import (
	"image/color"
	"testing"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// TODO: Test other direction than asecending

func TestValueWeightIntervalShouldCreate(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)
}

func TestValueWeightIntervalShouldSort(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	colors := []color.RGBA{
		{16, 16, 16, 255},
		{0, 0, 0, 255},
		{255, 255, 255, 255},
		{100, 100, 100, 255},
	}

	expectedResult := []color.Color{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{16, 16, 16, 255},
		color.RGBA{100, 100, 100, 255},
		color.RGBA{255, 255, 255, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortAscending)

	assert.Equal(t, expectedResult, actualResult)
}

func TestNormalizedWeightIntervalShouldCreate(t *testing.T) {
	interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
	assert.NotNil(t, interval)
}

func TestNormalizedWeightIntervalShouldSort(t *testing.T) {
	interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
	assert.NotNil(t, interval)

	colors := []color.RGBA{
		{16, 16, 16, 255},
		{0, 0, 0, 255},
		{255, 255, 255, 255},
		{100, 100, 100, 255},
	}

	expectedResult := []color.Color{
		color.RGBA{0, 0, 0, 255},
		color.RGBA{16, 16, 16, 255},
		color.RGBA{100, 100, 100, 255},
		color.RGBA{255, 255, 255, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortAscending)

	assert.Equal(t, expectedResult, actualResult)
}

// Create a test value weight determinant that is returning the red RGBA component as weight. Values from 0 to 255
func mockTestValueWeightDeterminant() func(color.RGBA) (int, error) {
	return func(c color.RGBA) (int, error) {
		r, _, _ := utils.RgbaToIntComponents(c)
		return r, nil
	}
}

// Create a test normalized weight determinant that is returning the red RGBA component as weight. Values from 0.0 to 1.0
func mockTestNormalizedWeightDeterminant() func(color.RGBA) (float64, error) {
	return func(c color.RGBA) (float64, error) {
		rNorm, _, _ := utils.RgbaToNormalizedComponents(c)
		return rNorm, nil
	}
}
