package sorter

import (
	"fmt"
	"image/color"
	"math/rand"
	"sort"
	"time"
)

const (
	defaultIntervalCapacity = 75
)

// Collection of image vertical or horizontal pixel neighbours with a propererty meeting some ceratin requirements
type Interval interface {
	// Add a RGBA color to the given interval
	Append(color color.RGBA) error

	// Get the count of colors stored in the interval
	Count() int

	// Get a boolean value representing the presence of elements in the given interval
	Any() bool

	// Sort all interval colors by weight in the specified direction and return the interval as a slice of RGBA colors
	Sort(direction SortDirection) []color.Color
}

// A interval interface implementation using a integer value for storing the interval element weight
type ValueWeightInterval struct {
	items                 []valueWeightIntervalItem
	weightDeterminantFunc func(color.RGBA) (int, error)
}

// Structure representing a single item in the value weight interval
type valueWeightIntervalItem struct {
	color  color.RGBA
	weight int
}

// A interval interface implementation using a floating point number value for storing the interval element weight
type NormalizedWeightInterval struct {
	items                 []normalizedWeightIntervalItem
	weightDeterminantFunc func(color.RGBA) (float64, error)
}

// Structure representing a single item in the normalized weight interval
type normalizedWeightIntervalItem struct {
	color  color.RGBA
	weight float64
}

func CreateValueWeightInterval(weightDeterminantFunc func(color.RGBA) (int, error)) Interval {
	interval := new(ValueWeightInterval)
	interval.items = make([]valueWeightIntervalItem, 0, defaultIntervalCapacity)
	interval.weightDeterminantFunc = weightDeterminantFunc

	return interval
}

func (interval *ValueWeightInterval) Append(color color.RGBA) error {
	weight, err := interval.weightDeterminantFunc(color)
	if err != nil {
		return fmt.Errorf("sorter: calculation of the color value weight failed: %w", err)
	}

	interval.items = append(interval.items, valueWeightIntervalItem{
		color:  color,
		weight: weight,
	})

	return nil
}

func (interval *ValueWeightInterval) Count() int {
	return len(interval.items)
}

func (interval *ValueWeightInterval) Any() bool {
	return interval.Count() > 0
}

func (interval *ValueWeightInterval) Sort(direction SortDirection) []color.Color {
	if len(interval.items) > 1 {
		switch direction {
		case SortAscending, SortDescending:
			{
				var sortDeterminantFunc func(i, j int) bool = nil

				if direction == SortAscending {
					sortDeterminantFunc = func(i, j int) bool {
						return interval.items[i].weight < interval.items[j].weight
					}
				}

				if direction == SortDescending {
					sortDeterminantFunc = func(i, j int) bool {
						return interval.items[i].weight > interval.items[j].weight
					}
				}

				if sortDeterminantFunc == nil {
					panic("sorter: undefined sort direction specified")
				}

				sort.Slice(interval.items, sortDeterminantFunc)
			}
		case SortRandom:
			{
				random := rand.New(rand.NewSource(time.Now().UnixNano()))
				random.Shuffle(len(interval.items), func(i, j int) {
					interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
				})
			}
		default:
			panic("sorter: undefined sort direction specified")
		}
	}

	intervalLength := len(interval.items)
	colors := make([]color.Color, intervalLength)

	for i := 0; i < intervalLength; i += 1 {
		colors[i] = interval.items[i].color
	}

	return colors
}

func CreateNormalizedWeightInterval(weightDeterminantFunc func(color.RGBA) (float64, error)) Interval {
	interval := new(NormalizedWeightInterval)
	interval.items = make([]normalizedWeightIntervalItem, 0, defaultIntervalCapacity)
	interval.weightDeterminantFunc = weightDeterminantFunc

	return interval
}

func (interval *NormalizedWeightInterval) Append(color color.RGBA) error {
	weight, err := interval.weightDeterminantFunc(color)
	if err != nil {
		return fmt.Errorf("sorter: calculation of the color normalized value weight failed: %w", err)
	}

	interval.items = append(interval.items, normalizedWeightIntervalItem{
		color:  color,
		weight: weight,
	})

	return nil
}

func (interval *NormalizedWeightInterval) Count() int {
	return len(interval.items)
}

func (interval *NormalizedWeightInterval) Any() bool {
	return interval.Count() > 0
}

func (interval *NormalizedWeightInterval) Sort(direction SortDirection) []color.Color {
	if len(interval.items) > 1 {
		switch direction {
		case SortAscending, SortDescending:
			{
				var sortDeterminantFunc func(i, j int) bool = nil

				if direction == SortAscending {
					sortDeterminantFunc = func(i, j int) bool {
						return interval.items[i].weight < interval.items[j].weight
					}
				}

				if direction == SortDescending {
					sortDeterminantFunc = func(i, j int) bool {
						return interval.items[i].weight > interval.items[j].weight
					}
				}

				if sortDeterminantFunc == nil {
					panic("sorter: undefined sort direction specified")
				}

				sort.Slice(interval.items, sortDeterminantFunc)
			}
		case SortRandom:
			{
				random := rand.New(rand.NewSource(time.Now().UnixNano()))
				random.Shuffle(len(interval.items), func(i, j int) {
					interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
				})
			}
		default:
			panic("sorter: undefined sort direction specified")
		}
	}

	intervalLength := len(interval.items)
	colors := make([]color.Color, intervalLength)

	for i := 0; i < intervalLength; i += 1 {
		colors[i] = interval.items[i].color
	}

	return colors
}
