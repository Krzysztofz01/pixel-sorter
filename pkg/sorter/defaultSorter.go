package sorter

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/sirupsen/logrus"
)

type defaultSorter struct {
	image   image.Image
	mask    *Mask
	options *SorterOptions
}

func CreateSorter(image image.Image, mask image.Image, options *SorterOptions) (Sorter, error) {
	sorter := new(defaultSorter)
	sorter.image = image

	if mask != nil {
		// TODO: The size validation can be moved to the mask factory func in the future
		if mask.Bounds().Dx() != image.Bounds().Dx() || mask.Bounds().Dy() != image.Bounds().Dy() {
			return nil, errors.New("sorter: the image and mask image sizes are not matching")
		}

		m, err := CreateMask(mask)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to create a new mask instance: %w", err)
		}

		sorter.mask = m
	} else {
		sorter.mask = CreateEmptyMask()
	}

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

	drawableImage = utils.RotateImage(drawableImage, sorter.options.Angle)

	switch sorter.options.SortOrder {
	case SortVertical:
		{
			if err := sorter.performParallelVerticalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the vertical sort")
			}
		}
	case SortHorizontal:
		{
			if err := sorter.performParallelHorizontalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the horizontal sort")
			}
		}
	case SortVerticalAndHorizontal:
		{
			if err := sorter.performParallelVerticalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the vertical sort")
			}

			if err := sorter.performParallelHorizontalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the horizontal sort")
			}
		}
	case SortHorizontalAndVertical:
		{
			if err := sorter.performParallelHorizontalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the horizontal sort")
			}

			if err := sorter.performParallelVerticalSort(&drawableImage); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the vertical sort")
			}
		}
	}

	drawableImage = utils.RotateImage(drawableImage, -sorter.options.Angle)
	drawableImage = utils.TrimImageTransparentWorkspace(drawableImage, sorter.image)

	return drawableImage, nil
}

func (sorter *defaultSorter) performHorizontalSort(drawableImage *draw.Image) error {
	for yIndex := 0; yIndex < (*drawableImage).Bounds().Dy(); yIndex += 1 {
		row, err := utils.GetImageRow(*drawableImage, yIndex)
		if err != nil {
			return fmt.Errorf("sorter: failed to retrieve the image pixel row for a given index: %w", err)
		}

		sortedRow, err := sorter.performSortOnImageStrip(row, func(iteratedCoordinate int) (int, int) {
			return iteratedCoordinate, yIndex
		})

		if err != nil {
			return fmt.Errorf("sorter: failed to perform the horizontal sorting: %w", err)
		}

		if err := utils.SetImageRow(drawableImage, sortedRow, yIndex); err != nil {
			return fmt.Errorf("sorter: failed to perform the insertion of the sorted row into the image: %w", err)
		}
	}

	return nil
}

// TODO: Check for potential race conditions
func (sorter *defaultSorter) performParallelHorizontalSort(drawableImage *draw.Image) error {
	yLength := (*drawableImage).Bounds().Dy()
	wg := sync.WaitGroup{}
	wg.Add(yLength)

	mu := sync.Mutex{}
	errCh := make(chan error)

	for y := 0; y < yLength; y += 1 {
		go func(yIndex int) {
			defer wg.Done()

			logrus.Debugf("Started to process row with index: %d", yIndex)

			row, err := utils.GetImageRow(*drawableImage, yIndex)
			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to retrieve the image pixel row for a given index: %w", err)
				return
			}

			sortedRow, err := sorter.performSortOnImageStrip(row, func(iteratedCoordinate int) (int, int) {
				return iteratedCoordinate, yIndex
			})

			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to perform the horizontal sorting: %w", err)
				return
			}

			mu.Lock()

			if err := utils.SetImageRow(drawableImage, sortedRow, yIndex); err != nil {
				errCh <- fmt.Errorf("sorter: failed to perform the insertion of the sorted row into the image: %w", err)
				mu.Unlock()
				return
			}

			mu.Unlock()

			logrus.Debugf("Finished to process row with index: %d", yIndex)
		}(y)
	}

	wg.Wait()
	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}

func (sorter *defaultSorter) performVerticalSort(drawableImage *draw.Image) error {
	for xIndex := 0; xIndex < (*drawableImage).Bounds().Dx(); xIndex += 1 {
		column, err := utils.GetImageColumn(*drawableImage, xIndex)
		if err != nil {
			return fmt.Errorf("sorter: failed to retrieve the image pixel column for a given index: %w", err)
		}

		sortedColumn, err := sorter.performSortOnImageStrip(column, func(iteratedCoordinate int) (int, int) {
			return xIndex, iteratedCoordinate
		})

		if err != nil {
			return fmt.Errorf("sorter: failed to perform the vertical sorting: %w", err)
		}

		if err := utils.SetImageColumn(drawableImage, sortedColumn, xIndex); err != nil {
			return fmt.Errorf("sorter: failed to perform the insertion of the sorted column into the image: %w", err)
		}
	}

	return nil
}

// TODO: Check for potential race conditions
func (sorter *defaultSorter) performParallelVerticalSort(drawableImage *draw.Image) error {
	xLength := (*drawableImage).Bounds().Dx()
	wg := sync.WaitGroup{}
	wg.Add(xLength)

	mu := sync.Mutex{}
	errCh := make(chan error)

	for x := 0; x < xLength; x += 1 {
		go func(xIndex int) {
			defer wg.Done()

			logrus.Debugf("Started to process column with index: %d", xIndex)

			column, err := utils.GetImageColumn(*drawableImage, xIndex)
			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to retrieve the image pixel column for a given index: %w", err)
				return
			}

			sortedColumn, err := sorter.performSortOnImageStrip(column, func(iteratedCoordinate int) (int, int) {
				return xIndex, iteratedCoordinate
			})

			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to perform the vertical sorting: %w", err)
				return
			}

			mu.Lock()

			if err := utils.SetImageColumn(drawableImage, sortedColumn, xIndex); err != nil {
				errCh <- fmt.Errorf("sorter: failed to perform the insertion of the sorted column into the image: %w", err)
				mu.Unlock()
				return
			}

			mu.Unlock()

			logrus.Debugf("Finished to process column with index: %d", xIndex)
		}(x)
	}

	wg.Wait()
	if len(errCh) > 0 {
		return <-errCh
	}

	return nil
}

// This is a helper function which performs the sorting of a given image strip by spliting it into intervals and sorting it by a given argument. This
// function is using the sorter preferences to determine how to perform the sort. There is also a maskedCoordinateFunc parameter which is a delegate
// used to retrieve information if a given pixel should be masked. We are using a external func for this in order to specify what coordinates should be
// looked up, beacuse this function has no access to the information which specific pixels from the image are processed now. Thanks to this approach, we
// can use a single function for both vertical and horizontal operations and just share a semi-fixed coordintes set.
func (sorter *defaultSorter) performSortOnImageStrip(imageStrip []color.Color, maskCoordinateFunc func(iteratedCoordinate int) (int, int)) ([]color.Color, error) {
	stripLength := len(imageStrip)
	sortedImageStrip := make([]color.Color, 0, stripLength)
	sortDirection := GetSortDeterminantDirection(sorter.options.SortDeterminant)

	interval := sorter.CreateInterval()
	for x := 0; x < stripLength; x += 1 {
		currentColor, err := utils.ColorToRgba(imageStrip[x])
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to convert the given color to a RGBA struct representation: %w", err)
		}

		isMasked, err := sorter.mask.IsMasked(maskCoordinateFunc(x))
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform a lookup to the mask image: %w", err)
		}

		// NOTE: isMasked and options dependecy solved using a quick K-Map
		passThrough := !isMasked || !sorter.options.UseMask

		if !utils.HasAnyTransparency(currentColor) && sorter.isMeetingIntervalRequirements(currentColor, isMasked) && passThrough {
			if err := interval.Append(currentColor); err != nil {
				return nil, fmt.Errorf("sorter: failed to append color to the interval: %w", err)
			}
		} else {
			if interval.Any() {
				sortedIntervalItems := interval.Sort(sortDirection)
				sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)

				interval = sorter.CreateInterval()
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

func (sorter *defaultSorter) isMeetingIntervalRequirements(color color.RGBA, isMasked bool) bool {
	tLower := sorter.options.IntervalDeterminantLowerThreshold
	tUpper := sorter.options.IntervalDeterminantUpperThreshold

	switch sorter.options.IntervalDeterminant {
	case SplitByBrightness:
		{
			brightness := utils.CalculatePerceivedBrightness(color)
			if brightness < tLower || brightness > tUpper {
				return false
			}

			return true
		}
	case SplitByHue:
		{
			h, _, _ := utils.RgbaToHsl(color)
			hNorm := float64(h) / 360.0

			if hNorm < tLower || hNorm > tUpper {
				return false
			}

			return true
		}
	case SplitByMask:
		{
			panic("sorter: not implemented")
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted interval determinant value")
	}
}

func (sorter *defaultSorter) CreateInterval() Interval {
	switch sorter.options.SortDeterminant {
	case SortByBrightnessAscending, SortByBrightnessDescending, ShuffleByBrightness:
		{
			return CreateNormalizedWeightInterval(func(c color.RGBA) (float64, error) {
				brightness := utils.CalculatePerceivedBrightness(c)
				return brightness, nil
			})
		}
	case SortByHueAscending, SortByHueDescending, ShuffleByHue:
		{
			return CreateValueWeightInterval(func(c color.RGBA) (int, error) {
				h, _, _ := utils.RgbaToHsl(c)
				return h, nil
			})
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted sorter weight determinant function value")
	}
}
