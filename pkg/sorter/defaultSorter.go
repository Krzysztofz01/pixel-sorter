package sorter

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

type defaultSorter struct {
	image   draw.Image
	options *SorterOptions
}

func CreateSorter(image image.Image, options *SorterOptions) (Sorter, error) {
	sorter := new(defaultSorter)

	drawableImage, err := utils.GetDrawableImage(image)
	if err != nil {
		return nil, fmt.Errorf("sorter: the provided image is not drawable: %w", err)
	}

	sorter.image = drawableImage

	if options != nil {
		lowerIdThreshold := options.IntervalDeterminantLowerThreshold
		if lowerIdThreshold > 1.0 || lowerIdThreshold < 0.0 {
			return nil, errors.New("sorter: invalid lower interval determinant threshold values provided")
		}

		upperIdThreshold := options.IntervalDeterminantUpperThreshold
		if upperIdThreshold > 1.0 || upperIdThreshold < 0.0 {
			return nil, errors.New("sorter: invalid upper interval determinant threshold values provided")
		}

		if lowerIdThreshold > upperIdThreshold {
			return nil, errors.New("sorter: the lower interval determiant threshold value can not be grate from the upper threshold")
		}

		sorter.options = options
	} else {
		sorter.options = GetDefaultSorterOptions()
	}

	return sorter, nil
}

func (sorter *defaultSorter) Sort() (image.Image, error) {
	xLength := sorter.image.Bounds().Dx()
	yLength := sorter.image.Bounds().Dy()

	// TODO: Rotate image

	for yIndex := 0; yIndex < yLength; yIndex += 1 {
		row, err := utils.GetImageRow(sorter.image, yIndex)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to retrieve the image pixel row for a given index: %w", err)
		}

		sortedRow, err := sorter.performSortOnImageStrip(row)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform the horizontal sorting: %w", err)
		}

		if err := utils.SetImageRow(&sorter.image, sortedRow, yIndex); err != nil {
			return nil, fmt.Errorf("sorter: failed to perform the insertion of the sorted row into the image: %w", err)
		}
	}

	for xIndex := 0; xIndex < xLength; xIndex += 1 {
		column, err := utils.GetImageColumn(sorter.image, xIndex)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to retrieve the image pixel column for a given index: %w", err)
		}

		sortedColumn, err := sorter.performSortOnImageStrip(column)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform the vertical sorting: %w", err)
		}

		if err := utils.SetImageColumn(&sorter.image, sortedColumn, xIndex); err != nil {
			return nil, fmt.Errorf("sorter: failed to perform the insertion of the sorted column into the image: %w", err)
		}
	}

	// TODO: Rotate image back

	return sorter.image, nil
}

func (sorter *defaultSorter) performSortOnImageStrip(imageStrip []color.Color) ([]color.Color, error) {
	stripLength := len(imageStrip)
	sortedImageStrip := make([]color.Color, 0)

	interval := CreateNormalizedWeightInterval(sorter.getWeightDeterminantFunction())
	for x := 0; x < stripLength; x += 1 {
		currentColor, err := utils.ColorToRgba(imageStrip[x])
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to convert the given color to a RGBA struct representation: %w", err)
		}

		if sorter.isMeetingIntervalRequirements(currentColor) {
			if err := interval.Append(currentColor); err != nil {
				return nil, fmt.Errorf("sorter: failed to append color to the interval: %w", err)
			}
		} else {
			if interval.Any() {
				sortedIntervalItems := interval.Sort()
				sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)

				interval = CreateNormalizedWeightInterval(sorter.getWeightDeterminantFunction())
			}

			sortedImageStrip = append(sortedImageStrip, currentColor)
		}
	}

	if interval.Any() {
		sortedIntervalItems := interval.Sort()
		sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)
	}

	return sortedImageStrip, nil
}

func (sorter *defaultSorter) isMeetingIntervalRequirements(color color.RGBA) bool {
	switch sorter.options.IntervalDeterminant {
	case SplitByBrightness:
		{
			tLower := sorter.options.IntervalDeterminantLowerThreshold
			tUpper := sorter.options.IntervalDeterminantUpperThreshold

			brightness := utils.CalculatePerceivedBrightness(color)
			if brightness < tLower || brightness > tUpper {
				return false
			}

			return true
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted interval determinant value")
	}
}

func (sorter *defaultSorter) getWeightDeterminantFunction() func(color.RGBA) (float64, error) {
	switch sorter.options.SortDeterminant {
	case SortByBrightness:
		{
			return func(c color.RGBA) (float64, error) {
				brightness := utils.CalculatePerceivedBrightness(c)
				return brightness, nil
			}
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted sorter weight determinant function value")
	}
}
