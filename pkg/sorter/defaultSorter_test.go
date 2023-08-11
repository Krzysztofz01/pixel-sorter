package sorter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

// TODO: Currently the test are verifying that no errors occure but the resulting image is not verified

func TestDefaultOptionsAndSortDeterminantBrightness(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByBrightness

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDeterminantHue(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByHue

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDeterminantSaturation(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortBySaturation

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDeterminantRedChannel(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByRedChannel

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDeterminantGreenChannel(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByGreenChannel

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDeterminantBlueChannel(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDeterminant = SortByBlueChannel

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndAngle45Degrees(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Angle = 45

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndCycles3(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Cycles = 3

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDirectionAscending(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = SortAscending

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDirectionDescending(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = SortDescending

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDirectionShuffle(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = Shuffle

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndSortDirectionRandom(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortDirection = SortRandom

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalDeterminantBrightness(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByBrightness

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalDeterminantHue(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByHue

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalDeterminantSaturation(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitBySaturation

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalDeterminantMask(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByMask

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), mockTestBlackAndWhiteStripesImage(), nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalDeterminantAbsolute(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByAbsoluteColor

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalDeterminantEdge(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminant = SplitByEdgeDetection

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndLowerIntervalThreshold04UpperIntervalThreshold06(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalDeterminantLowerThreshold = 0.4
	options.IntervalDeterminantUpperThreshold = 0.6

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalMaxLength2(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalLength = 2

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalMaxLength1RandomFactor1(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.IntervalLength = 1
	options.IntervalLengthRandomFactor = 1

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndUseMask(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.UseMask = true

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), mockTestBlackAndWhiteStripesImage(), nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndOrderHorizontal(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortHorizontal

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndOrderVertical(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortVertical

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndOrderHorizontalVertical(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortHorizontalAndVertical

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndOrderVerticalHorizontal(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.SortOrder = SortVerticalAndHorizontal

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndScale05(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Scale = 0.5

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndBlendingModeNone(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Blending = BlendingNone

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndBlendingModeLighten(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Blending = BlendingLighten

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndBlendingModeDarken(t *testing.T) {
	defer goleak.VerifyNone(t)

	options := GetDefaultSorterOptions()
	options.Blending = BlendingDarken

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalPaintingFill(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalFill

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalPaintingGradient(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalGradient

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalPaintingRepeat(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalRepeat

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

func TestDefaultOptionsAndIntervalPaintingAverage(t *testing.T) {
	options := GetDefaultSorterOptions()
	options.IntervalPainting = IntervalAverage

	sorter, err := CreateSorter(mockTestBlackAndWhiteStripesImage(), nil, nil, options)

	assert.NotNil(t, sorter)
	assert.Nil(t, err)

	result, err := sorter.Sort()

	assert.NotNil(t, result)
	assert.Nil(t, err)
}

const (
	mock_image_width  = 5
	mock_image_height = 5
)

// Create a test image which consists of black and white 1px wide columns of the size specified by the mock_image prefixed constants
func mockTestBlackAndWhiteStripesImage() draw.Image {
	p1 := image.Point{0, 0}
	p2 := image.Point{mock_image_width, mock_image_height}
	image := image.NewRGBA(image.Rectangle{p1, p2})

	cBlack := color.RGBA{0, 0, 0, 0xff}
	cWhite := color.RGBA{255, 255, 255, 0xff}

	for xIndex := 0; xIndex < mock_image_width; xIndex += 1 {
		for yIndex := 0; yIndex < mock_image_height; yIndex += 1 {
			if xIndex%2 != 0 {
				image.Set(xIndex, yIndex, cBlack)
			} else {
				image.Set(xIndex, yIndex, cWhite)
			}
		}
	}

	return image
}
