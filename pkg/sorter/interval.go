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

// Collection of image vertical or horizontal pixel neighbours with a propererty meeting some ceratin requirements
type Interval interface {
	// Add a RGBA color to the given interval
	Append(color color.RGBA) error

	// Get the count of colors stored in the interval
	Count() int

	// Get a boolean value representing the presence of elements in the given interval
	Any() bool

	// Sort all interval colors by weight in the specified direction and return the interval as a new slice of RGBA
	// colors. The internal interval items collection will be cleared after the sort.
	Sort(direction SortDirection, painting IntervalPainting) []color.RGBA

	// Sort all interval colors by weight in the specified direction and return the interval by writing the RGBA
	// colors to the provided buffer. The internal interval items collection will be cleared after the sort.
	SortToBuffer(direction SortDirection, painting IntervalPainting, buffer *[]color.RGBA)
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
				h, _, _, _ := utils.RgbaToHsla(c)
				return h
			})
		}
	case SortBySaturation:
		{
			return CreateNormalizedWeightInterval(func(c color.RGBA) float64 {
				_, s, _, _ := utils.RgbaToHsla(c)
				return s
			})
		}
	case SortByAbsoluteColor:
		{
			return CreateValueWeightInterval(func(c color.RGBA) int {
				return int(c.R) * int(c.G) * int(c.B)
			})
		}
	case SortByRedChannel:
		{
			return CreateValueWeightInterval(func(c color.RGBA) int {
				return int(c.R)
			})
		}
	case SortByGreenChannel:
		{
			return CreateValueWeightInterval(func(c color.RGBA) int {
				return int(c.G)
			})
		}
	case SortByBlueChannel:
		{
			return CreateValueWeightInterval(func(c color.RGBA) int {
				return int(c.B)
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

func (interval *genericInterval[T]) Sort(direction SortDirection, painting IntervalPainting) []color.RGBA {
	buffer := make([]color.RGBA, 0, interval.Count())

	interval.SortToBuffer(direction, painting, &buffer)

	return buffer
}

func (interval *genericInterval[T]) SortToBuffer(direction SortDirection, painting IntervalPainting, buffer *[]color.RGBA) {
	defer func() {
		// TODO: The previous items are not garbage-collected after the "clear" operation and can lead to pseudo memory leaks.
		interval.items = interval.items[:0]
	}()

	if interval.Count() <= 1 {
		for i := 0; i < interval.Count(); i += 1 {
			*buffer = append(*buffer, interval.items[i].color)
		}

		return
	}

	switch painting {
	case IntervalRepeat:
		{
			for range interval.items {
				*buffer = append(*buffer, interval.items[0].color)
			}
			return
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

			for range interval.items {
				*buffer = append(*buffer, avg)
			}

			return
		}
	case IntervalFill:
		{
			if direction == SortRandom {
				if utils.CIntn(2) == 1 {
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
					// TODO: Implement a math/rand independent shuffle
					random := rand.New(rand.NewSource(time.Now().UnixNano()))
					random.Shuffle(len(interval.items), func(i, j int) {
						interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
					})
				}
			default:
				panic("sorter: undefined sort direction specified")
			}

			for i := 0; i < interval.Count(); i += 1 {
				*buffer = append(*buffer, interval.items[i].color)
			}

			return
		}
	case IntervalGradient:
		{
			if direction == SortRandom {
				if utils.CIntn(2) == 1 {
					direction = SortAscending
				} else {
					direction = SortDescending
				}
			}

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

						abColorLerp := utils.InterpolateRgba(a.color, b.color, t)
						bcColorLerp := utils.InterpolateRgba(b.color, c.color, t)

						*buffer = append(*buffer, utils.InterpolateRgba(abColorLerp, bcColorLerp, t))
					}
				}
			case Shuffle:
				{
					// TODO: Implement a meth/rand independent shuffle
					random := rand.New(rand.NewSource(time.Now().UnixNano()))
					random.Shuffle(len(interval.items), func(i, j int) {
						interval.items[i], interval.items[j] = interval.items[j], interval.items[i]
					})

					a := interval.items[0].color
					b := interval.items[(interval.Count()-1)/2].color
					c := interval.items[interval.Count()-1].color

					for i := 0; i < interval.Count(); i += 1 {
						t := float64(i) / float64(interval.Count()-1)

						abColorLerp := utils.InterpolateRgba(a, b, t)
						bcColorLerp := utils.InterpolateRgba(b, c, t)

						*buffer = append(*buffer, utils.InterpolateRgba(abColorLerp, bcColorLerp, t))
					}
				}
			default:
				panic("sorter: undefined sort direction specified")
			}
		}
	default:
		panic("sorter: undefined interval painting specified")
	}
}
