package sorter

import (
	"image/color"
	"sort"
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

func TestValueWeightIntervalShouldPaintRepeat(t *testing.T) {
	cases := []SortDirection{
		SortAscending,
		SortDescending,
		Shuffle,
		SortRandom,
	}

	for _, sortDirection := range cases {
		interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
		assert.NotNil(t, interval)

		colors := []color.RGBA{
			{16, 16, 16, 255},
			{0, 0, 0, 255},
			{255, 255, 255, 255},
			{100, 100, 100, 255},
		}

		expectedResult := []color.Color{
			color.RGBA{16, 16, 16, 255},
			color.RGBA{16, 16, 16, 255},
			color.RGBA{16, 16, 16, 255},
			color.RGBA{16, 16, 16, 255},
		}

		assert.False(t, interval.Any())

		for _, color := range colors {
			err := interval.Append(color)
			assert.Nil(t, err)
		}

		assert.True(t, interval.Any())

		actualResult := interval.Sort(sortDirection, IntervalRepeat)

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestValueWeightIntervalShouldPaintAverage(t *testing.T) {
	cases := []SortDirection{
		SortAscending,
		SortDescending,
		Shuffle,
		SortRandom,
	}

	for _, sortDirection := range cases {
		interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
		assert.NotNil(t, interval)

		colors := []color.RGBA{
			{16, 16, 16, 255},
			{0, 0, 0, 255},
			{255, 255, 255, 255},
			{100, 100, 100, 255},
		}

		expectedResult := []color.Color{
			color.RGBA{92, 92, 92, 255},
			color.RGBA{92, 92, 92, 255},
			color.RGBA{92, 92, 92, 255},
			color.RGBA{92, 92, 92, 255},
		}

		assert.False(t, interval.Any())

		for _, color := range colors {
			err := interval.Append(color)
			assert.Nil(t, err)
		}

		assert.True(t, interval.Any())

		actualResult := interval.Sort(sortDirection, IntervalAverage)

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestValueWeightIntervalShouldSortAscendingPaintFill(t *testing.T) {
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

	actualResult := interval.Sort(SortAscending, IntervalFill)

	assert.Equal(t, expectedResult, actualResult)
}

func TestValueWeightIntervalShouldSortDescendingPaintFill(t *testing.T) {
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

	actualResult := interval.Sort(SortDescending, IntervalFill)

	assert.Equal(t, expectedResult, actualResult)
}

func TestValueWeightIntervalShouldShufflePaintFill(t *testing.T) {
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

	actualResult := interval.Sort(Shuffle, IntervalFill)

	assert.ElementsMatch(t, colors, actualResult)
}

func TestValueWeightIntervalShouldSortRandomPaintFill(t *testing.T) {
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

	sortedColors := interval.Sort(SortRandom, IntervalFill)

	isSortedAscending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R > right.R
	})

	assert.False(t, !isSortedAscending && !isSortedDescending)
	assert.ElementsMatch(t, colors, sortedColors)
}

func TestValueWeightIntervalShouldSortAscendingPaintGradient(t *testing.T) {
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
		color.RGBA{72, 72, 72, 255},
		color.RGBA{157, 157, 157, 255},
		color.RGBA{255, 255, 255, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortAscending, IntervalGradient)

	assert.Equal(t, expectedResult, actualResult)
}

func TestValueWeightIntervalShouldSortDescendingPaintGradient(t *testing.T) {
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
		color.RGBA{157, 157, 157, 255},
		color.RGBA{72, 72, 72, 255},
		color.RGBA{0, 0, 0, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortDescending, IntervalGradient)

	assert.Equal(t, expectedResult, actualResult)
}

func TestValueWeightIntervalShouldShufflePaintGradient(t *testing.T) {
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

	interval.Sort(Shuffle, IntervalGradient)

	// TODO: Implement first and last elements assertion
}

func TestValueWeightIntervalShouldSortRandomPaintGradient(t *testing.T) {
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

	sortedColors := interval.Sort(SortRandom, IntervalGradient)

	isSortedAscending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R > right.R
	})

	assert.False(t, !isSortedAscending && !isSortedDescending)
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

func TestNormalizedWeightIntervalShouldPaintRepeat(t *testing.T) {
	cases := []SortDirection{
		SortAscending,
		SortDescending,
		Shuffle,
		SortRandom,
	}

	for _, sortDirection := range cases {
		interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
		assert.NotNil(t, interval)

		colors := []color.RGBA{
			{16, 16, 16, 255},
			{0, 0, 0, 255},
			{255, 255, 255, 255},
			{100, 100, 100, 255},
		}

		expectedResult := []color.Color{
			color.RGBA{16, 16, 16, 255},
			color.RGBA{16, 16, 16, 255},
			color.RGBA{16, 16, 16, 255},
			color.RGBA{16, 16, 16, 255},
		}

		assert.False(t, interval.Any())

		for _, color := range colors {
			err := interval.Append(color)
			assert.Nil(t, err)
		}

		assert.True(t, interval.Any())

		actualResult := interval.Sort(sortDirection, IntervalRepeat)

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestNormalizedWeightIntervalShouldPaintAverage(t *testing.T) {
	cases := []SortDirection{
		SortAscending,
		SortDescending,
		Shuffle,
		SortRandom,
	}

	for _, sortDirection := range cases {
		interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
		assert.NotNil(t, interval)

		colors := []color.RGBA{
			{16, 16, 16, 255},
			{0, 0, 0, 255},
			{255, 255, 255, 255},
			{100, 100, 100, 255},
		}

		expectedResult := []color.Color{
			color.RGBA{92, 92, 92, 255},
			color.RGBA{92, 92, 92, 255},
			color.RGBA{92, 92, 92, 255},
			color.RGBA{92, 92, 92, 255},
		}

		assert.False(t, interval.Any())

		for _, color := range colors {
			err := interval.Append(color)
			assert.Nil(t, err)
		}

		assert.True(t, interval.Any())

		actualResult := interval.Sort(sortDirection, IntervalAverage)

		assert.Equal(t, expectedResult, actualResult)
	}
}

func TestNormalizedWeightIntervalShouldSortAscendingPaintFill(t *testing.T) {
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

	actualResult := interval.Sort(SortAscending, IntervalFill)

	assert.Equal(t, expectedResult, actualResult)
}

func TestNormalizedWeightIntervalShouldSortDescendingPaintFill(t *testing.T) {
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

	actualResult := interval.Sort(SortDescending, IntervalFill)

	assert.Equal(t, expectedResult, actualResult)
}

func TestNormalizedWeightIntervalShouldShufflePaintFill(t *testing.T) {
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

	actualResult := interval.Sort(Shuffle, IntervalFill)

	assert.ElementsMatch(t, colors, actualResult)
}

func TestNormalizedWeightIntervalShouldSortRandomPaintFill(t *testing.T) {
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

	sortedColors := interval.Sort(SortRandom, IntervalFill)

	isSortedAscending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R > right.R
	})

	assert.False(t, !isSortedAscending && !isSortedDescending)
	assert.ElementsMatch(t, colors, sortedColors)
}

func TestNormalizedWeightIntervalShouldSortAscendingPaintGradient(t *testing.T) {
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
		color.RGBA{72, 72, 72, 255},
		color.RGBA{157, 157, 157, 255},
		color.RGBA{255, 255, 255, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortAscending, IntervalGradient)

	assert.Equal(t, expectedResult, actualResult)
}

func TestNormalizedWeightIntervalShouldSortDescendingPaintGradient(t *testing.T) {
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
		color.RGBA{157, 157, 157, 255},
		color.RGBA{72, 72, 72, 255},
		color.RGBA{0, 0, 0, 255},
	}

	assert.False(t, interval.Any())

	for _, color := range colors {
		err := interval.Append(color)
		assert.Nil(t, err)
	}

	assert.True(t, interval.Any())

	actualResult := interval.Sort(SortDescending, IntervalGradient)

	assert.Equal(t, expectedResult, actualResult)
}

func TestNormalizedWeightIntervalShouldShufflePaintGradient(t *testing.T) {
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

	interval.Sort(Shuffle, IntervalGradient)

	// TODO: Implement first and last elements assertion
}

func TestNormalizedWeightIntervalShouldSortRandomPaintGradient(t *testing.T) {
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

	sortedColors := interval.Sort(SortRandom, IntervalGradient)

	isSortedAscending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left, _ := sortedColors[i].(color.RGBA)
		right, _ := sortedColors[j].(color.RGBA)

		return left.R > right.R
	})

	assert.False(t, !isSortedAscending && !isSortedDescending)
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
