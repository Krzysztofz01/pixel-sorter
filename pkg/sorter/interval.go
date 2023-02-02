package sorter

import (
	"fmt"
	"image/color"
	"sort"
)

// Collection of image vertical or horizontal pixel neighbours with a propererty meeting some ceratin requirements
type Interval interface {
	// Add a RGBA color to the given interval
	Append(color color.RGBA) error

	// Get a boolean value representing the presence of elements in the given interval
	Any() bool

	// Sort all interval colors by weight and return the interval as a slice of RGBA colors
	Sort() []color.Color
}

// A interval interface implementation using a integer value for storing the interval element weight
type ValueWeightInterval struct {
	items             []valueWeightIntervalItem
	weightDeterminant func(color.RGBA) (int, error)
}

// Structure representing a single item in the value weight interval
type valueWeightIntervalItem struct {
	color  color.RGBA
	weight int
}

// A interval interface implementation using a floating point number value for storing the interval element weight
type NormalizedWeightInterval struct {
	items             []normalizedWeightIntervalItem
	weightDeterminant func(color.RGBA) (float64, error)
}

// Structure representing a single item in the normalized weight interval
type normalizedWeightIntervalItem struct {
	color  color.RGBA
	weight float64
}

func CreateValueWeightInterval(weightDeterminant func(color.RGBA) (int, error)) Interval {
	interval := new(ValueWeightInterval)
	// TODO: Fine-tune ths size
	interval.items = make([]valueWeightIntervalItem, 0)
	interval.weightDeterminant = weightDeterminant

	return interval
}

func (interval *ValueWeightInterval) Append(color color.RGBA) error {
	weight, err := interval.weightDeterminant(color)
	if err != nil {
		return fmt.Errorf("sorter: calculation of the color value weight failed: %w", err)
	}

	interval.items = append(interval.items, valueWeightIntervalItem{
		color:  color,
		weight: weight,
	})

	return nil
}

func (interval *ValueWeightInterval) Any() bool {
	return len(interval.items) > 0
}

func (interval *ValueWeightInterval) Sort() []color.Color {
	sort.Slice(interval.items, func(i, j int) bool {
		return interval.items[i].weight < interval.items[j].weight
	})

	intervalLength := len(interval.items)
	colors := make([]color.Color, intervalLength)

	for i := 0; i < intervalLength; i += 1 {
		colors[i] = interval.items[i].color
	}

	return colors
}

func CreateNormalizedWeightInterval(weightDeterminant func(color.RGBA) (float64, error)) Interval {
	interval := new(NormalizedWeightInterval)
	// TODO: Fine-tune ths size
	interval.items = make([]normalizedWeightIntervalItem, 0)
	interval.weightDeterminant = weightDeterminant

	return interval
}

func (interval *NormalizedWeightInterval) Append(color color.RGBA) error {
	weight, err := interval.weightDeterminant(color)
	if err != nil {
		return fmt.Errorf("sorter: calculation of the color normalized value weight failed: %w", err)
	}

	interval.items = append(interval.items, normalizedWeightIntervalItem{
		color:  color,
		weight: weight,
	})

	return nil
}

func (interval *NormalizedWeightInterval) Any() bool {
	return len(interval.items) > 0
}

func (interval *NormalizedWeightInterval) Sort() []color.Color {
	sort.Slice(interval.items, func(i, j int) bool {
		return interval.items[i].weight < interval.items[j].weight
	})

	intervalLength := len(interval.items)
	colors := make([]color.Color, intervalLength)

	for i := 0; i < intervalLength; i += 1 {
		colors[i] = interval.items[i].color
	}

	return colors
}
