package sorter

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

const (
	defaultIntervalCapacity = 75
)

// TODO: After the Sort() the internal items collection order is changed, and some sort algorithms are assuming that
// the "current" order is the append order. We can either create a local copy of the items collection, restore the
// internal collection after the sort or prohibit the interval to sort more than one time. Currently there is no case
// of multiple sorts on a single interval.

// Collection of image vertical or horizontal pixel neighbours with a propererty meeting some ceratin requirements
type Interval interface {
	// Add a RGBA color to the given interval
	Append(color color.RGBA) error

	// Get the count of colors stored in the interval
	Count() int

	// Get a boolean value representing the presence of elements in the given interval
	Any() bool

	// Sort all interval colors by weight in the specified direction and return the interval as a slice of RGBA colors
	Sort(direction SortDirection, painting IntervalPainting) []color.Color
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

func (interval *ValueWeightInterval) Sort(direction SortDirection, painting IntervalPainting) []color.Color {
	if interval.Count() < 1 {
		colors := make([]color.Color, interval.Count())
		for i := 0; i < interval.Count(); i += 1 {
			colors[i] = interval.items[i].color
		}

		return colors
	}

	switch painting {
	case IntervalRepeat:
		{
			colors := make([]color.Color, interval.Count())
			for i := 0; i < interval.Count(); i += 1 {
				colors[i] = interval.items[0].color
			}

			return colors
		}
	case IntervalAverage:
		{
			sumR, sumG, sumB, sumA := 0, 0, 0, 0
			for _, item := range interval.items {
				sumR += int(item.color.R)
				sumG += int(item.color.G)
				sumB += int(item.color.B)
				sumA += int(item.color.A)
			}

			avg := color.RGBA{
				R: uint8(sumR / interval.Count()),
				G: uint8(sumG / interval.Count()),
				B: uint8(sumB / interval.Count()),
				A: uint8(sumA / interval.Count()),
			}

			colors := make([]color.Color, interval.Count())
			for i := 0; i < interval.Count(); i += 1 {
				colors[i] = avg
			}

			return colors
		}
	case IntervalFill:
		{
			if direction == SortRandom {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))

				if random.Intn(2) == 1 {
					direction = SortAscending
				} else {
					direction = SortDescending
				}
			}

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
			case Shuffle:
				{
					random := rand.New(rand.NewSource(time.Now().UnixNano()))
					random.Shuffle(len(interval.items), func(i, j int) {
						interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
					})
				}
			default:
				panic("sorter: undefined sort direction specified")
			}

			colors := make([]color.Color, interval.Count())
			for i := 0; i < interval.Count(); i += 1 {
				colors[i] = interval.items[i].color
			}

			return colors
		}
	case IntervalGradient:
		{
			if direction == SortRandom {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))

				if random.Intn(2) == 1 {
					direction = SortAscending
				} else {
					direction = SortDescending
				}
			}

			colors := make([]color.Color, interval.Count())

			switch direction {
			case SortAscending, SortDescending:
				{
					a := interval.items[0]
					c := interval.items[0]
					for _, item := range interval.items {
						if direction == SortAscending {
							if item.weight < a.weight {
								a = item
							}

							if item.weight > c.weight {
								c = item
							}

							continue
						}

						if direction == SortDescending {
							if item.weight > a.weight {
								a = item
							}

							if item.weight < c.weight {
								c = item
							}

							continue
						}

						panic("sorter: invalid sort direction state for gradient painting")
					}

					b := interval.items[0]
					bLerp := utils.Lerp(float64(a.weight), float64(c.weight), 0.5)
					bDelta := math.Inf(1)

					for _, item := range interval.items {
						delta := math.Abs(float64(item.weight) - bLerp)
						if delta < bDelta {
							bDelta = delta
							b = item
						}
					}

					for i := 0; i < interval.Count(); i += 1 {
						t := float64(i) / float64(interval.Count()-1)

						abColorLerp := utils.InterpolateColor(a.color, b.color, t)
						bcColorLerp := utils.InterpolateColor(b.color, c.color, t)

						colors[i] = utils.InterpolateColor(abColorLerp, bcColorLerp, t)
					}
				}
			case Shuffle:
				{
					random := rand.New(rand.NewSource(time.Now().UnixNano()))
					random.Shuffle(len(interval.items), func(i, j int) {
						interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
					})

					a := interval.items[0].color
					b := interval.items[(interval.Count()-1)/2].color
					c := interval.items[interval.Count()-1].color

					for i := 0; i < interval.Count(); i += 1 {
						t := float64(i) / float64(interval.Count()-1)

						abColorLerp := utils.InterpolateColor(a, b, t)
						bcColorLerp := utils.InterpolateColor(b, c, t)

						colors[i] = utils.InterpolateColor(abColorLerp, bcColorLerp, t)
					}
				}
			default:
				panic("sorter: undefined sort direction specified")
			}

			return colors
		}
	default:
		panic("sorter: undefined interval painting specified")
	}
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

func (interval *NormalizedWeightInterval) Sort(direction SortDirection, painting IntervalPainting) []color.Color {
	if interval.Count() < 1 {
		colors := make([]color.Color, interval.Count())
		for i := 0; i < interval.Count(); i += 1 {
			colors[i] = interval.items[i].color
		}

		return colors
	}

	switch painting {
	case IntervalRepeat:
		{
			colors := make([]color.Color, interval.Count())
			for i := 0; i < interval.Count(); i += 1 {
				colors[i] = interval.items[0].color
			}

			return colors
		}
	case IntervalAverage:
		{
			sumR, sumG, sumB, sumA := 0, 0, 0, 0
			for _, item := range interval.items {
				sumR += int(item.color.R)
				sumG += int(item.color.G)
				sumB += int(item.color.B)
				sumA += int(item.color.A)
			}

			avg := color.RGBA{
				R: uint8(sumR / interval.Count()),
				G: uint8(sumG / interval.Count()),
				B: uint8(sumB / interval.Count()),
				A: uint8(sumA / interval.Count()),
			}

			colors := make([]color.Color, interval.Count())
			for i := 0; i < interval.Count(); i += 1 {
				colors[i] = avg
			}

			return colors
		}
	case IntervalFill:
		{
			if direction == SortRandom {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))

				if random.Intn(2) == 1 {
					direction = SortAscending
				} else {
					direction = SortDescending
				}
			}

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
			case Shuffle:
				{
					random := rand.New(rand.NewSource(time.Now().UnixNano()))
					random.Shuffle(len(interval.items), func(i, j int) {
						interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
					})
				}
			default:
				panic("sorter: undefined sort direction specified")
			}

			colors := make([]color.Color, interval.Count())
			for i := 0; i < interval.Count(); i += 1 {
				colors[i] = interval.items[i].color
			}

			return colors
		}
	case IntervalGradient:
		{
			if direction == SortRandom {
				random := rand.New(rand.NewSource(time.Now().UnixNano()))

				if random.Intn(2) == 1 {
					direction = SortAscending
				} else {
					direction = SortDescending
				}
			}

			colors := make([]color.Color, interval.Count())

			switch direction {
			case SortAscending, SortDescending:
				{
					a := interval.items[0]
					c := interval.items[0]
					for _, item := range interval.items {
						if direction == SortAscending {
							if item.weight < a.weight {
								a = item
							}

							if item.weight > c.weight {
								c = item
							}

							continue
						}

						if direction == SortDescending {
							if item.weight > a.weight {
								a = item
							}

							if item.weight < c.weight {
								c = item
							}

							continue
						}

						panic("sorter: invalid sort direction state for gradient painting")
					}

					b := interval.items[0]
					bLerp := utils.Lerp(a.weight, c.weight, 0.5)
					bDelta := math.Inf(1)

					for _, item := range interval.items {
						delta := math.Abs(item.weight - bLerp)
						if delta < bDelta {
							bDelta = delta
							b = item
						}
					}

					for i := 0; i < interval.Count(); i += 1 {
						t := float64(i) / float64(interval.Count()-1)

						abColorLerp := utils.InterpolateColor(a.color, b.color, t)
						bcColorLerp := utils.InterpolateColor(b.color, c.color, t)

						colors[i] = utils.InterpolateColor(abColorLerp, bcColorLerp, t)
					}
				}
			case Shuffle:
				{
					random := rand.New(rand.NewSource(time.Now().UnixNano()))
					random.Shuffle(len(interval.items), func(i, j int) {
						interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
					})

					a := interval.items[0].color
					b := interval.items[(interval.Count()-1)/2].color
					c := interval.items[interval.Count()-1].color

					for i := 0; i < interval.Count(); i += 1 {
						t := float64(i) / float64(interval.Count()-1)

						abColorLerp := utils.InterpolateColor(a, b, t)
						bcColorLerp := utils.InterpolateColor(b, c, t)

						colors[i] = utils.InterpolateColor(abColorLerp, bcColorLerp, t)
					}
				}
			default:
				panic("sorter: undefined sort direction specified")
			}

			return colors
		}
	default:
		panic("sorter: undefined interval painting specified")
	}
}
