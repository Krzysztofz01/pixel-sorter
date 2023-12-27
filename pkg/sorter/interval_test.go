package sorter

import (
	"image/color"
	"sort"
	"testing"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateIntervalShouldCreateIntervalForSortByBrightness(t *testing.T) {
	interval := CreateInterval(SortByBrightness)
	assert.NotNil(t, interval)
}

func TestCreateIntervalShouldCreateIntervalForSortByHue(t *testing.T) {
	interval := CreateInterval(SortByHue)
	assert.NotNil(t, interval)
}

func TestCreateIntervalShouldCreateIntervalForSortBySaturation(t *testing.T) {
	interval := CreateInterval(SortBySaturation)
	assert.NotNil(t, interval)
}

func TestCreateIntervalShouldCreateIntervalForSortByAbsolute(t *testing.T) {
	interval := CreateInterval(SortByAbsoluteColor)
	assert.NotNil(t, interval)
}

func TestCreateIntervalShouldCreateIntervalForSortByRedChannel(t *testing.T) {
	interval := CreateInterval(SortByRedChannel)
	assert.NotNil(t, interval)
}

func TestCreateIntervalShouldCreateIntervalForSortByGreenChannel(t *testing.T) {
	interval := CreateInterval(SortByGreenChannel)
	assert.NotNil(t, interval)
}

func TestCreateIntervalShouldCreateIntervalForSortByBlueChannel(t *testing.T) {
	interval := CreateInterval(SortByBlueChannel)
	assert.NotNil(t, interval)
}

func TestCreateIntervalShouldPanicForInvalidSortDeterminant(t *testing.T) {
	assert.Panics(t, func() {
		CreateInterval(-1)
	})
}

func TestValueWeightIntervalShouldCreate(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)
}

func TestValueWeightIntervalShouldNotCreateForNilWeightDeterminantFunc(t *testing.T) {
	assert.Panics(t, func() {
		CreateValueWeightInterval(nil)
	})
}

func TestNormalizedIntervalShouldNotCreateForNilWeightDeterminantFunc(t *testing.T) {
	assert.Panics(t, func() {
		CreateNormalizedWeightInterval(nil)
	})
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

func TestValueWeightIntervalShouldSortWhenContainingSingleElement(t *testing.T) {
	interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
	assert.NotNil(t, interval)

	color := color.RGBA{0, 0, 0, 255}
	err := interval.Append(color)
	assert.Nil(t, err)

	result := interval.Sort(SortAscending, IntervalFill)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
	assert.Equal(t, color, result[0])
}

func TestValueWeightIntervalShouldPanicForInvalidDirectionAndPaintingValues(t *testing.T) {
	cases := []struct {
		d SortDirection
		p IntervalPainting
	}{
		{-100, IntervalFill},
		{-100, IntervalGradient},
		{-100, -100},
	}

	for _, c := range cases {
		assert.Panics(t, func() {
			interval := CreateValueWeightInterval(mockTestValueWeightDeterminant())
			assert.NotNil(t, interval)

			err := interval.Append(color.RGBA{0, 0, 0, 255})
			assert.Nil(t, err)

			err = interval.Append(color.RGBA{255, 255, 255, 255})
			assert.Nil(t, err)

			interval.Sort(c.d, c.p)
		})
	}
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

		expectedResult := []color.RGBA{
			{16, 16, 16, 255},
			{16, 16, 16, 255},
			{16, 16, 16, 255},
			{16, 16, 16, 255},
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

		expectedResult := []color.RGBA{
			{92, 92, 92, 255},
			{92, 92, 92, 255},
			{92, 92, 92, 255},
			{92, 92, 92, 255},
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

	expectedResult := []color.RGBA{
		{0, 0, 0, 255},
		{16, 16, 16, 255},
		{100, 100, 100, 255},
		{255, 255, 255, 255},
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

	expectedResult := []color.RGBA{
		{255, 255, 255, 255},
		{100, 100, 100, 255},
		{16, 16, 16, 255},
		{0, 0, 0, 255},
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
		left := sortedColors[i]
		right := sortedColors[j]

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left := sortedColors[i]
		right := sortedColors[j]

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

	expectedResult := []color.RGBA{
		{0, 0, 0, 255},
		{72, 72, 72, 255},
		{157, 157, 157, 255},
		{255, 255, 255, 255},
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

	expectedResult := []color.RGBA{
		{255, 255, 255, 255},
		{157, 157, 157, 255},
		{72, 72, 72, 255},
		{0, 0, 0, 255},
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
		left := sortedColors[i]
		right := sortedColors[j]

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left := sortedColors[i]
		right := sortedColors[j]

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

func TestNormalziedWeightIntervalShouldSortWhenContainingSingleElement(t *testing.T) {
	interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
	assert.NotNil(t, interval)

	color := color.RGBA{0, 0, 0, 255}
	err := interval.Append(color)
	assert.Nil(t, err)

	result := interval.Sort(SortAscending, IntervalFill)

	assert.NotNil(t, result)
	assert.NotEmpty(t, result)
	assert.Equal(t, color, result[0])
}

func TestNormalizedWeightIntervalShouldPanicForInvalidDirectionAndPaintingValues(t *testing.T) {
	cases := []struct {
		d SortDirection
		p IntervalPainting
	}{
		{-100, IntervalFill},
		{-100, IntervalGradient},
		{-100, -100},
	}

	for _, c := range cases {
		assert.Panics(t, func() {
			interval := CreateNormalizedWeightInterval(mockTestNormalizedWeightDeterminant())
			assert.NotNil(t, interval)

			err := interval.Append(color.RGBA{0, 0, 0, 255})
			assert.Nil(t, err)

			err = interval.Append(color.RGBA{255, 255, 255, 255})
			assert.Nil(t, err)

			interval.Sort(c.d, c.p)
		})
	}
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

		expectedResult := []color.RGBA{
			{16, 16, 16, 255},
			{16, 16, 16, 255},
			{16, 16, 16, 255},
			{16, 16, 16, 255},
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

		expectedResult := []color.RGBA{
			{92, 92, 92, 255},
			{92, 92, 92, 255},
			{92, 92, 92, 255},
			{92, 92, 92, 255},
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

	expectedResult := []color.RGBA{
		{0, 0, 0, 255},
		{16, 16, 16, 255},
		{100, 100, 100, 255},
		{255, 255, 255, 255},
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

	expectedResult := []color.RGBA{
		{255, 255, 255, 255},
		{100, 100, 100, 255},
		{16, 16, 16, 255},
		{0, 0, 0, 255},
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
		left := sortedColors[i]
		right := sortedColors[j]

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left := sortedColors[i]
		right := sortedColors[j]

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

	expectedResult := []color.RGBA{
		{0, 0, 0, 255},
		{72, 72, 72, 255},
		{157, 157, 157, 255},
		{255, 255, 255, 255},
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

	expectedResult := []color.RGBA{
		{255, 255, 255, 255},
		{157, 157, 157, 255},
		{72, 72, 72, 255},
		{0, 0, 0, 255},
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
		left := sortedColors[i]
		right := sortedColors[j]

		return left.R < right.R
	})

	isSortedDescending := sort.SliceIsSorted(sortedColors, func(i, j int) bool {
		left := sortedColors[i]
		right := sortedColors[j]

		return left.R > right.R
	})

	assert.False(t, !isSortedAscending && !isSortedDescending)
}

// Create a test value weight determinant that is returning the red RGBA component as weight. Values from 0 to 255
func mockTestValueWeightDeterminant() func(color.RGBA) int {
	return func(c color.RGBA) int {
		return int(c.R)
	}
}

// Create a test normalized weight determinant that is returning the red RGBA component as weight. Values from 0.0 to 1.0
func mockTestNormalizedWeightDeterminant() func(color.RGBA) float64 {
	return func(c color.RGBA) float64 {
		rNorm, _, _ := utils.RgbaToNormalizedComponents(c)
		return rNorm
	}
}
