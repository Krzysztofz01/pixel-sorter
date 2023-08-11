package sorter

import (
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

type genericInterval[T int | float64] struct {
	items                 []genericIntervalItem[T]
	weightDeterminantFunc func(color.RGBA) T
}

type genericIntervalItem[T int | float64] struct {
	color  color.RGBA
	weight T
}

// Create a new interval instance with the item weights represented as a integer values
func CreateValueWeightInterval(weightDeterminantFunc func(color.RGBA) int) Interval {
	if weightDeterminantFunc == nil {
		panic("sorter: the provided weight determinant function is nil")
	}

	return &genericInterval[int]{
		items:                 make([]genericIntervalItem[int], 0, defaultIntervalCapacity),
		weightDeterminantFunc: weightDeterminantFunc,
	}
}

// Create a new interval instance with the item weights represented as normalzied values
func CreateNormalizedWeightInterval(weightDeterminantFunc func(color.RGBA) float64) Interval {
	if weightDeterminantFunc == nil {
		panic("sorter: the provided weight determinant function is nil")
	}

	return &genericInterval[float64]{
		items:                 make([]genericIntervalItem[float64], 0, defaultIntervalCapacity),
		weightDeterminantFunc: weightDeterminantFunc,
	}
}

// Create a new interval instance based on the specifications required by the provided sort determinant
func CreateInterval(sort SortDeterminant) Interval {
	switch sort {
	case SortByBrightness:
		{
			return CreateNormalizedWeightInterval(func(c color.RGBA) float64 {
				return utils.CalculatePerceivedBrightness(c)
			})
		}
	case SortByHue:
		{
			return CreateValueWeightInterval(func(c color.RGBA) int {
				h, _, _, _ := utils.ColorToHsla(c)
				return h
			})
		}
	case SortBySaturation:
		{
			return CreateNormalizedWeightInterval(func(c color.RGBA) float64 {
				_, s, _, _ := utils.ColorToHsla(c)
				return s
			})
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted sorter weight determinant function value")
	}
}

func (interval *genericInterval[T]) Append(color color.RGBA) error {
	weight := interval.weightDeterminantFunc(color)
	interval.items = append(interval.items, genericIntervalItem[T]{
		color:  color,
		weight: weight,
	})

	return nil
}

func (interval *genericInterval[T]) Count() int {
	return len(interval.items)
}

func (interval *genericInterval[T]) Any() bool {
	return len(interval.items) > 0
}

func (interval *genericInterval[T]) Sort(direction SortDirection, painting IntervalPainting) []color.Color {
	if interval.Count() <= 1 {
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
			case SortAscending:
				{
					sort.SliceStable(interval.items, func(i, j int) bool {
						return interval.items[i].weight < interval.items[j].weight
					})
				}
			case SortDescending:
				{
					sort.SliceStable(interval.items, func(i, j int) bool {
						return interval.items[i].weight > interval.items[j].weight
					})
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
						// NOTE: According to the case statement values, assuming that the direction is not ascending it must be descending
						if direction == SortAscending {
							if item.weight < a.weight {
								a = item
							}

							if item.weight > c.weight {
								c = item
							}
						} else {
							if item.weight > a.weight {
								a = item
							}

							if item.weight < c.weight {
								c = item
							}
						}
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
