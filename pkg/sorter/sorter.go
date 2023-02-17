package sorter

import (
	"image"
)

// Flag representing the determinant parameter for the sorting process
type SortDeterminant int

// TODO: Make the direction as a separate parameter
const (
	SortByBrightnessAscending SortDeterminant = iota
	SortByBrightnessDescending
	ShuffleByBrightness
	SortByHueAscending
	SortByHueDescending
	ShuffleByHue
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
)

// Structure representing all the parameters for the sorter
type SorterOptions struct {
	SortDeterminant                   SortDeterminant
	SortOrder                         SortOrder
	IntervalDeterminant               IntervalDeterminant
	IntervalDeterminantLowerThreshold float64
	IntervalDeterminantUpperThreshold float64
	IntervalLength                    int
	Angle                             int
	UseMask                           bool
	Cycles                            int
}

// Get a SorterOptions structure instance with default values
func GetDefaultSorterOptions() *SorterOptions {
	options := new(SorterOptions)
	options.Angle = 0
	options.SortDeterminant = SortByBrightnessAscending
	options.SortOrder = SortHorizontalAndVertical
	options.IntervalDeterminant = SplitByBrightness
	options.IntervalDeterminantLowerThreshold = 0.0
	options.IntervalDeterminantUpperThreshold = 1.0
	options.UseMask = false
	options.IntervalLength = 0
	options.Cycles = 1

	return options
}

// Utility used to create a pixel sorted version of a given image
type Sorter interface {
	// Perform the sorting operation and return the sorted version of the image
	Sort() (image.Image, error)
}

func GetSortDeterminantDirection(s SortDeterminant) SortDirection {
	switch s {
	case SortByBrightnessAscending, SortByHueAscending:
		{
			return SortAscending
		}
	case SortByBrightnessDescending, SortByHueDescending:
		{
			return SortDescending
		}
	case ShuffleByBrightness, ShuffleByHue:
		{
			return SortRandom
		}
	default:
		panic("sorter: invalid sort determinant specified")
	}
}
