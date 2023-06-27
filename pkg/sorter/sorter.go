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
	Shuffle
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
	IntervalLengthRandomFactor        int
	Angle                             int
	UseMask                           bool
	Cycles                            int
	Scale                             float64
	Blending                          ResultImageBlending
}

// Return a boolean value indicating if the given sorter options combination is valid
// and a string containing validation failure message if the options came out to be invalid
func (options *SorterOptions) AreValid() (bool, string) {
	if options.IntervalDeterminantLowerThreshold < 0.0 || options.IntervalDeterminantLowerThreshold > 1.0 {
		return false, "lower interval determinant threshold must be between values 0 and 1"
	}

	if options.IntervalDeterminantUpperThreshold < 0.0 || options.IntervalDeterminantUpperThreshold > 1.0 {
		return false, "upper interval determinant threshold must be between values 0 and 1"
	}

	if options.IntervalDeterminantLowerThreshold > options.IntervalDeterminantUpperThreshold {
		return false, "lower interval determinant threshold must no be greater than the upper one"
	}

	if options.Cycles < 1 {
		return false, "the cycles count must be 1 or greater"
	}

	if options.Scale <= 0.0 || options.Scale > 1.0 {
		return false, "the scale factor must be between values 0 (exclusive) and 1"
	}

	if options.IntervalLength < 0 {
		return false, "the interval max length values must not be negative"
	}

	if options.IntervalLengthRandomFactor < 0 {
		return false, "the interval length random factor value must not be negative"
	}

	return true, ""
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
	options.IntervalLengthRandomFactor = 0
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

// Utility used to create a pixel sorted version of a given video
type VideoSorter interface {
	// Perform the sorting operation and store the sorted version in the sorter specified path
	Sort() error
}
