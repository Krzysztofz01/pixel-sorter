package sorter

import (
	"image"
)

// Flag representing the determinant parameter for the sorting process
type SortDeterminant int

const (
	SortByBrightness SortDeterminant = iota
	SortByHue
	SortBySaturation
)

// Flag representing the order in which should be the image sorted
type SortOrder int

const (
	SortHorizontal SortOrder = iota
	SortVertical
	SortHorizontalAndVertical
	SortVerticalAndHorizontal
)

// Flag representing the direction of the sorting
type SortDirection int

const (
	SortAscending SortDirection = iota
	SortDescending
	// NOTE: The naming here is ugly... shuffling is not sorting
	SortRandom
)

// Flag representing the determinant for spliting the image into intervals
type IntervalDeterminant int

const (
	SplitByBrightness IntervalDeterminant = iota
	SplitByHue
	SplitBySaturation
	SplitByMask
	SplitByAbsoluteColor
	SplitByEdgeDetection
)

type ResultImageBlending int

const (
	BlendingNone ResultImageBlending = iota
	BlendingLighten
	BlendingDarken
)

// Structure representing all the parameters for the sorter
type SorterOptions struct {
	SortDeterminant                   SortDeterminant
	SortDirection                     SortDirection
	SortOrder                         SortOrder
	IntervalDeterminant               IntervalDeterminant
	IntervalDeterminantLowerThreshold float64
	IntervalDeterminantUpperThreshold float64
	IntervalLength                    int
	Angle                             int
	UseMask                           bool
	Cycles                            int
	Scale                             float64
	Blending                          ResultImageBlending
}

// Get a SorterOptions structure instance with default values
func GetDefaultSorterOptions() *SorterOptions {
	options := new(SorterOptions)
	options.Angle = 0
	options.SortDeterminant = SortByBrightness
	options.SortDirection = SortAscending
	options.SortOrder = SortHorizontalAndVertical
	options.IntervalDeterminant = SplitByBrightness
	options.IntervalDeterminantLowerThreshold = 0.0
	options.IntervalDeterminantUpperThreshold = 1.0
	options.UseMask = false
	options.IntervalLength = 0
	options.Cycles = 1
	options.Scale = 1
	options.Blending = BlendingNone

	return options
}

// Utility used to create a pixel sorted version of a given image
type Sorter interface {
	// Perform the sorting operation and return the sorted version of the image
	Sort() (image.Image, error)
}
