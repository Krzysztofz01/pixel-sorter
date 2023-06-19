package sorter

import (
	"image/color"
	"testing"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestValueWeightIntervalShouldCreate(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)
}

func TestValueWeightIntervalShouldTellIfContainsAnyColors(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	assert.False(t, interval.Any())

	err := interval.Append(color.RGBA{0, 0, 0, 255})
	assert.Nil(t, err)

	assert.True(t, interval.Any())
}

func TestValueWeightIntervalShouldTellTheCountOfContainedColors(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	assert.Equal(t, 0, interval.Count())

	err := interval.Append(color.RGBA{0, 0, 0, 255})
	assert.Nil(t, err)

	assert.Equal(t, 1, interval.Count())
}

func TestValueWeightIntervalShouldSortAscending(t *testing.T) {
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

func TestValueWeightIntervalShouldSortDescending(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	colors := []color.RGBA{
		{16, 16, 16, 255},
		{0, 0, 0, 255},
		{255, 255, 255, 255},
		{100, 100, 100, 255},
	}

	expectedResult := []color.Color{
		color.RGBA{255, 255, 255, 255},
		color.RGBA{100, 100, 100, 255},
		color.RGBA{16, 16, 16, 255},
		color.RGBA{0, 0, 0, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortDescending)

	assert.Equal(t, expectedResult, actualResult)
}

func TestValueWeightIntervalShouldShuffle(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	colors := []color.RGBA{
		{16, 16, 16, 255},
		{0, 0, 0, 255},
		{255, 255, 255, 255},
		{100, 100, 100, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(Shuffle)

	assert.ElementsMatch(t, colors, actualResult)
}

func TestNormalizedWeightIntervalShouldCreate(t *testing.T) {
	interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
	assert.NotNil(t, interval)
}

func TestNormalizedWeightIntervalShouldTellIfContainsAnyColors(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	assert.False(t, interval.Any())

	err := interval.Append(color.RGBA{0, 0, 0, 255})
	assert.Nil(t, err)

	assert.True(t, interval.Any())
}

func TestNormalizedWeightIntervalShouldTellTheCountOfContainedColors(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	assert.Equal(t, 0, interval.Count())

	err := interval.Append(color.RGBA{0, 0, 0, 255})
	assert.Nil(t, err)

	assert.Equal(t, 1, interval.Count())
}

func TestNormalizedWeightIntervalShouldSortAscending(t *testing.T) {
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

func TestNormalizedWeightIntervalShouldSortDescending(t *testing.T) {
	interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
	assert.NotNil(t, interval)

	colors := []color.RGBA{
		{16, 16, 16, 255},
		{0, 0, 0, 255},
		{255, 255, 255, 255},
		{100, 100, 100, 255},
	}

	expectedResult := []color.Color{
		color.RGBA{255, 255, 255, 255},
		color.RGBA{100, 100, 100, 255},
		color.RGBA{16, 16, 16, 255},
		color.RGBA{0, 0, 0, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortDescending)

	assert.Equal(t, expectedResult, actualResult)
}

func TestNormalizedWeightIntervalShouldShuffle(t *testing.T) {
	interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
	assert.NotNil(t, interval)

	colors := []color.RGBA{
		{16, 16, 16, 255},
		{0, 0, 0, 255},
		{255, 255, 255, 255},
		{100, 100, 100, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(Shuffle)

	assert.ElementsMatch(t, colors, actualResult)
}

// Create a test value weight determinant that is returning the red RGBA component as weight. Values from 0 to 255
func mockTestValueWeightDeterminant() func(color.RGBA) (int, error) {
	return func(c color.RGBA) (int, error) {
		return int(c.R), nil
	}
}

// Create a test normalized weight determinant that is returning the red RGBA component as weight. Values from 0.0 to 1.0
func mockTestNormalizedWeightDeterminant() func(color.RGBA) (float64, error) {
	return func(c color.RGBA) (float64, error) {
		rNorm, _, _ := utils.RgbaToNormalizedComponents(c)
		return rNorm, nil
	}
}
