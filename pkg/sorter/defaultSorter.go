package sorter

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/img"
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
			return nil, errors.New("sorter: the lower interval determiant threshold value can not be greater from the upper threshold")
		}

		if options.Cycles < 1 {
			return nil, errors.New("sorter: the cycles count can not be zero or less")
		}

		if options.Scale < 0.0 || options.Scale > 1.0 {
			return nil, errors.New("sorter: the scale percentage must be in range between zero and one")
		}

		sorter.options = options
	} else {
		sorter.options = GetDefaultSorterOptions()
	}

	if sorter.options.Scale != 1.0 {
		imageScaled, err := utils.ScaleImage(sorter.image, sorter.options.Scale)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to scale the target image: %w", err)
		}

		sorter.image = imageScaled
	}

	if sorter.options.IntervalDeterminant == SplitByEdgeDetection {
		edgeDetectionExecTime := time.Now()
		imageEdges, err := img.PerformEdgeDetection(sorter.image, false)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform the edge detection on the provided image: %w", err)
		}

		invertedEdges, err := utils.InvertImage(imageEdges)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform color inversion on the edge detection image: %w", err)
		}

		mask = invertedEdges
		logrus.Debugf("Edge detection took: %s.", time.Since(edgeDetectionExecTime))
	}

	if mask != nil {
		maskExecTime := time.Now()

		if sorter.options.Scale != 1.0 {
			scaledMask, err := utils.ScaleImage(mask, sorter.options.Scale)
			if err != nil {
				return nil, fmt.Errorf("sorter: failed to scale the target image mask: %w", err)
			}

			mask = scaledMask
		}

		m, err := CreateImageMask(mask, image.Bounds())
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to create a new mask instance: %w", err)
		}

		logrus.Debugf("Mask parsing took: %s.", time.Since(maskExecTime))
		sorter.mask = m
	} else {
		sorter.mask = CreateEmptyMask()
	}

	return sorter, nil
}

func (sorter *defaultSorter) Sort() (image.Image, error) {
	sortingExecTime := time.Now()
	drawableImage, err := utils.GetDrawableImage(sorter.image)
	if err != nil {
		return nil, fmt.Errorf("sorter: the provided image is not drawable: %w", err)
	}

	drawableImage = utils.RotateImage(drawableImage, sorter.options.Angle)

	for c := 0; c < sorter.options.Cycles; c += 1 {
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
	}

	drawableImage = utils.RotateImage(drawableImage, -sorter.options.Angle)
	drawableImage = utils.TrimImageTransparentWorkspace(drawableImage, sorter.image)

	switch sorter.options.Blending {
	case BlendingLighten:
		{
			if drawableImage, err = utils.BlendImages(sorter.image, drawableImage, utils.LightenOnly); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the image blending: %w", err)
			}
		}
	case BlendingDarken:
		{
			if drawableImage, err = utils.BlendImages(sorter.image, drawableImage, utils.DarkenOnly); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the image blending: %w", err)
			}
		}
	case BlendingNone:
		break
	default:
		panic("sorter: invalid blending mode specified")
	}

	logrus.Debugf("Pixel sorting took: %s.", time.Since(sortingExecTime))
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

func (sorter *defaultSorter) performParallelHorizontalSort(drawableImage *draw.Image) error {
	yLength := (*drawableImage).Bounds().Dy()
	wg := sync.WaitGroup{}
	wg.Add(yLength)

	mu := sync.RWMutex{}
	errCh := make(chan error)

	for y := 0; y < yLength; y += 1 {
		go func(yIndex int) {
			defer wg.Done()

			mu.RLock()

			row, err := utils.GetImageRow(*drawableImage, yIndex)
			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to retrieve the image pixel row for a given index: %w", err)
				mu.RUnlock()
				return
			}

			mu.RUnlock()

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

func (sorter *defaultSorter) performParallelVerticalSort(drawableImage *draw.Image) error {
	xLength := (*drawableImage).Bounds().Dx()
	wg := sync.WaitGroup{}
	wg.Add(xLength)

	mu := sync.RWMutex{}
	errCh := make(chan error)

	for x := 0; x < xLength; x += 1 {
		go func(xIndex int) {
			defer wg.Done()

			mu.RLock()

			column, err := utils.GetImageColumn(*drawableImage, xIndex)
			if err != nil {
				errCh <- fmt.Errorf("sorter: failed to retrieve the image pixel column for a given index: %w", err)
				mu.RUnlock()
				return
			}

			mu.RUnlock()

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

		if !utils.HasAnyTransparency(currentColor) && sorter.isMeetingIntervalRequirements(currentColor, isMasked, interval) && passThrough {
			if err := interval.Append(currentColor); err != nil {
				return nil, fmt.Errorf("sorter: failed to append color to the interval: %w", err)
			}
		} else {
			if interval.Any() {
				sortedIntervalItems := interval.Sort(sorter.options.SortDirection)
				sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)

				interval = sorter.CreateInterval()
			}

			sortedImageStrip = append(sortedImageStrip, currentColor)
		}
	}

	if interval.Any() {
		sortedIntervalItems := interval.Sort(sorter.options.SortDirection)
		sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)
	}

	return sortedImageStrip, nil
}

func (sorter *defaultSorter) isMeetingIntervalRequirements(color color.RGBA, isMasked bool, interval Interval) bool {
	// NOTE: interval length and options dependecy solved using a quick K-Map
	maxLength := sorter.options.IntervalLength
	if !(maxLength == 0) && (maxLength <= interval.Count()) {
		return false
	}

	switch sorter.options.IntervalDeterminant {
	case SplitByBrightness:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			brightness := utils.CalculatePerceivedBrightness(color)
			return brightness >= lThreshold && brightness <= uThreshold
		}
	case SplitByHue:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			h, _, _ := utils.RgbaToHsl(color)
			hNorm := float64(h) / 360.0

			return hNorm >= lThreshold && hNorm <= uThreshold
		}
	case SplitBySaturation:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			_, s, _ := utils.RgbaToHsl(color)
			return s >= lThreshold && s <= uThreshold
		}
	case SplitByMask, SplitByEdgeDetection:
		{
			return !isMasked
		}
	case SplitByAbsoluteColor:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			r, g, b := utils.RgbaToIntComponents(color)
			abs := float64((r * g * b)) / 16581375.0

			return abs >= lThreshold && abs < uThreshold
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted interval determinant value")
	}
}

func (sorter *defaultSorter) CreateInterval() Interval {
	switch sorter.options.SortDeterminant {
	case SortByBrightness:
		{
			return CreateNormalizedWeightInterval(func(c color.RGBA) (float64, error) {
				brightness := utils.CalculatePerceivedBrightness(c)
				return brightness, nil
			})
		}
	case SortByHue:
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
