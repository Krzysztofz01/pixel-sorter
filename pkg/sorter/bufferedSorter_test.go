package sorter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

// TODO: Currently the test are verifying that no errors occure but the resulting image is not verified

func TestBufferedSorterSortingCancellationShouldBreakTheSorting(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	sorter, err := CreateBufferedSorter(mockHugeTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	sortingGoroutine := func() {
		result, err := sorter.Sort(options)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrSortingCancellation)
	}

	cancellationGoroutine := func() {
		assert.True(t, sorter.CancelSort())

		assert.False(t, sorter.CancelSort())
	}

	go sortingGoroutine()
	go cancellationGoroutine()
}

func TestBufferedSorterDefaultOptionsAndSortDeterminantBrightness(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByBrightness

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDeterminantHue(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByHue

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDeterminantSaturation(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortBySaturation

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDeterminantAbsoluteColor(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByAbsoluteColor

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDeterminantRedChannel(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByRedChannel

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDeterminantGreenChannel(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByGreenChannel

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDeterminantBlueChannel(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByBlueChannel

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndAngle45Degrees(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Angle = 45

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndCycles3(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Cycles = 3

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDirectionAscending(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = SortAscending

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDirectionDescending(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = SortDescending

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDirectionShuffle(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = Shuffle

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndSortDirectionRandom(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = SortRandom

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalDeterminantBrightness(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByBrightness

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalDeterminantHue(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByHue

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalDeterminantSaturation(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitBySaturation

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalDeterminantMask(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByMask

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), mockTestBlackAndWhiteStripesImage(), nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalDeterminantAbsolute(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByAbsoluteColor

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalDeterminantEdge(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByEdgeDetection

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndLowerIntervalThreshold04UpperIntervalThreshold06(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminantLowerThreshold = 0.4
	options.IntervalDeterminantUpperThreshold = 0.6

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalMaxLength2(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalLength = 2

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalMaxLength1RandomFactor1(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalLength = 1
	options.IntervalLengthRandomFactor = 1

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndUseMask(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.UseMask = true

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), mockTestBlackAndWhiteStripesImage(), nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndOrderHorizontal(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortHorizontal

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndOrderVertical(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortVertical

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndOrderHorizontalVertical(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortHorizontalAndVertical

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndOrderVerticalHorizontal(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortVerticalAndHorizontal

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndScale05(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Scale = 0.5

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndBlendingModeNone(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Blending = BlendingNone

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndBlendingModeLighten(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Blending = BlendingLighten

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndBlendingModeDarken(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Blending = BlendingDarken

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalPaintingFill(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalFill

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalPaintingGradient(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalGradient

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalPaintingRepeat(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalRepeat

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestBufferedSorterDefaultOptionsAndIntervalPaintingAverage(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalAverage

	sorter, err := CreateBufferedSorter(mockTestBlackAndWhiteStripesImage(), nil, nil)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort(options)
	assert.NotNil(t, result)
	assert.Nil(t, err)
	result, err = sorter.Sort(options)

	assert.NotNil(t, result)
	assert.Nil(t, err)
}
