package sorter

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

type defaultSorter struct {
	image   image.Image
	options *SorterOptions
}

func CreateSorter(image image.Image, options *SorterOptions) (Sorter, error) {
	sorter := new(defaultSorter)
	sorter.image = image

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
	drawableImage, err := utils.GetDrawableImage(sorter.image)
	if err != nil {
		return nil, fmt.Errorf("sorter: the provided image is not drawable: %w", err)
	}

	// TODO: Implement image rotation here
	// TODO: Add support for trimming the transparent pixels during the sorting process
	// TODO: Add support for handing situations where the xLength and yLength is changing due to the rotation

	switch sorter.options.SortOrder {
	case SortVertical:
		{
			if err := sorter.performVerticalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the vertical sort")
			}
		}
	case SortHorizontal:
		{
			if err := sorter.performHorizontalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the horizontal sort")
			}
		}
	case SortVerticalAndHorizontal:
		{
			if err := sorter.performVerticalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the vertical sort")
			}

			if err := sorter.performHorizontalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the horizontal sort")
			}
		}
	case SortHorizontalAndVertical:
		{
			if err := sorter.performHorizontalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the horizontal sort")
			}

			if err := sorter.performVerticalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the vertical sort")
			}
		}
	}

	// TODO: Implement image back rotation here
	// TODO: Implement the workspace trimming here

	return drawableImage, nil
}

func (sorter *defaultSorter) performHorizontalSort(drawableImage *draw.Image) error {
	// TODO: Can be this value cached somehow?
	yLength := (*drawableImage).Bounds().Dy()

	for yIndex := 0; yIndex < yLength; yIndex += 1 {
		row, err := utils.GetImageRow(*drawableImage, yIndex)
		if err != nil {
			return fmt.Errorf("sorter: failed to retrieve the image pixel row for a given index: %w", err)
		}

		sortedRow, err := sorter.performSortOnImageStrip(row)
		if err != nil {
			return fmt.Errorf("sorter: failed to perform the horizontal sorting: %w", err)
		}

		if err := utils.SetImageRow(drawableImage, sortedRow, yIndex); err != nil {
			return fmt.Errorf("sorter: failed to perform the insertion of the sorted row into the image: %w", err)
		}
	}

	return nil
}

// TODO: First time using go rutines and sync. code... More research is required!!!
func (sorter *defaultSorter) performParallelHorizontalSort(drawableImage *draw.Image) error {
	yLength := (*drawableImage).Bounds().Dy()

	//mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(yLength)

	errCh := make(chan error)

	for yIndex := 0; yIndex < yLength; yIndex += 1 {
		go func(y int) {
			defer wg.Done()

			row, err := utils.GetImageRow(*drawableImage, y)
			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to retrieve the image pixel row for a given index: %w", err)
				return
			}

			sortedRow, err := sorter.performSortOnImageStrip(row)
			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to perform the horizontal sorting: %w", err)
				return
			}

			//mu.Lock()
			if err := utils.SetImageRow(drawableImage, sortedRow, y); err != nil {
				errCh <- fmt.Errorf("sorter: failed to perform the insertion of the sorted row into the image: %w", err)
				//mu.Unlock()
				return
			}
			//mu.Unlock()
		}(yIndex)
	}

	wg.Wait()
	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}

func (sorter *defaultSorter) performVerticalSort(drawableImage *draw.Image) error {
	// TODO: Can be this value cached somehow?
	xLength := (*drawableImage).Bounds().Dx()

	for xIndex := 0; xIndex < xLength; xIndex += 1 {
		column, err := utils.GetImageColumn(*drawableImage, xIndex)
		if err != nil {
			return fmt.Errorf("sorter: failed to retrieve the image pixel column for a given index: %w", err)
		}

		sortedColumn, err := sorter.performSortOnImageStrip(column)
		if err != nil {
			return fmt.Errorf("sorter: failed to perform the vertical sorting: %w", err)
		}

		if err := utils.SetImageColumn(drawableImage, sortedColumn, xIndex); err != nil {
			return fmt.Errorf("sorter: failed to perform the insertion of the sorted column into the image: %w", err)
		}
	}

	return nil
}

func (sorter *defaultSorter) performSortOnImageStrip(imageStrip []color.Color) ([]color.Color, error) {
	stripLength := len(imageStrip)
	sortedImageStrip := make([]color.Color, 0, stripLength)
	sortDirection := GetSortDeterminantDirection(sorter.options.SortDeterminant)

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
				sortedIntervalItems := interval.Sort(sortDirection)
				sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)

				interval = CreateNormalizedWeightInterval(sorter.getWeightDeterminantFunction())
			}

			sortedImageStrip = append(sortedImageStrip, currentColor)
		}
	}

	if interval.Any() {
		sortedIntervalItems := interval.Sort(sortDirection)
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

// TODO: Implement support for hue weight
func (sorter *defaultSorter) getWeightDeterminantFunction() func(color.RGBA) (float64, error) {
	switch sorter.options.SortDeterminant {
	case SortByBrightnessAscending, SortByBrightnessDescending, ShuffleByBrightness:
		{
			return func(c color.RGBA) (float64, error) {
				brightness := utils.CalculatePerceivedBrightness(c)
				return brightness, nil
			}
		}
	case SortByHueAscending, SortByHueDescending, ShuffleByHue:
		{
			panic("sorter: not implemented")
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted sorter weight determinant function value")
	}
}
