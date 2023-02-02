package sorter

import (
	"image"
)

// Flag representing the determinant parameter for the sorting process
type SortDeterminant string

const (
	SortByBrightness SortDeterminant = "brightness"
	SortByHue        SortDeterminant = "hue"
)

// Flag representing the determinant for spliting the image into intervals
type IntervalDeterminant string

const (
	SplitByBrightness IntervalDeterminant = "brightness"
	//SplitByHue        IntervalDeterminant = "hue"
)

// Structure representing all the parameters for the sorter
type SorterOptions struct {
	SortDeterminant                   SortDeterminant
	IntervalDeterminant               IntervalDeterminant
	IntervalDeterminantLowerThreshold float64
	IntervalDeterminantUpperThreshold float64
	Angle                             int
}

// Get a SorterOptions structure instance with default values
func GetDefaultSorterOptions() *SorterOptions {
	options := new(SorterOptions)
	options.Angle = 0
	options.SortDeterminant = SortByBrightness
	options.IntervalDeterminant = SplitByBrightness
	// TODO: Fine-tune this
	options.IntervalDeterminantLowerThreshold = 0.0
	options.IntervalDeterminantUpperThreshold = 1.0

	return options
}

// Utility used to create a pixel sorted version of a given image
type Sorter interface {
	// Perform the sorting operation and return the sorted version of the image
	Sort() (image.Image, error)
}
