package sorter

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"sync"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/img"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/sirupsen/logrus"
)

type defaultSorter struct {
	image   image.Image
	mask    *Mask
	logger  *logrus.Entry
	options *SorterOptions
}

// Create a new image sorter instance by providing the image to be sorted and optional parameters such as mask image
// logger instance and custom sorter options. This function will return a new sorter instance or a error.
func CreateSorter(image image.Image, mask image.Image, logger *logrus.Logger, options *SorterOptions) (Sorter, error) {
	sorter := new(defaultSorter)
	sorter.image = image

	if logger == nil {
		loggerBuffer := bytes.Buffer{}
		defaultLogger := &logrus.Logger{
			Out:       &loggerBuffer,
			Formatter: &logrus.TextFormatter{},
		}

		sorter.logger = defaultLogger.WithField("prefix", "pixel-sorter")
	} else {
		sorter.logger = logger.WithField("prefix", "pixel-sorter")
	}

	if options != nil {
		sorter.logger.Debugln("Running the sorter with specified sorter options.")

		if valid, msg := options.AreValid(); !valid {
			sorter.logger.Debugf("Sorter options validation failed. Sorter options: %+v", *options)
			return nil, fmt.Errorf("sorter: %s", msg)
		}

		sorter.logger.Debugf("Sorter options validation passed. Sorter options: %+v", *options)
		sorter.options = options
	} else {
		sorter.logger.Debugln("Running the sorter with default sorter options.")
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
		sorter.logger.Debugf("Edge detection took: %s.", time.Since(edgeDetectionExecTime))
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

		m, err := CreateImageMask(mask, sorter.image.Bounds(), sorter.options.Angle)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to create a new mask instance: %w", err)
		}

		sorter.logger.Debugf("Mask parsing took: %s.", time.Since(maskExecTime))
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

	if sorter.options.Angle != 0 {
		drawableImage = utils.RotateImage(drawableImage, sorter.options.Angle)
	}

	for c := 0; c < sorter.options.Cycles; c += 1 {
		switch sorter.options.SortOrder {
		case SortVertical:
			{
				if err := sorter.performParallelVerticalSort(drawableImage); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical sort: %w", err)
				}
			}
		case SortHorizontal:
			{
				if err := sorter.performParallelHorizontalSort(drawableImage); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal sort: %w", err)
				}
			}
		case SortVerticalAndHorizontal:
			{
				if err := sorter.performParallelVerticalSort(drawableImage); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical sort: %w", err)
				}

				if err := sorter.performParallelHorizontalSort(drawableImage); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal sort: %w", err)
				}
			}
		case SortHorizontalAndVertical:
			{
				if err := sorter.performParallelHorizontalSort(drawableImage); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal sort: %w", err)
				}

				if err := sorter.performParallelVerticalSort(drawableImage); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical sort: %w", err)
				}
			}
		}
	}

	if sorter.options.Angle != 0 {
		drawableImage = utils.RotateImage(drawableImage, -sorter.options.Angle)
		drawableImage = utils.TrimImageTransparentWorkspace(drawableImage, sorter.image)
	}

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

	sorter.logger.Debugf("Pixel sorting took: %s.", time.Since(sortingExecTime))
	return drawableImage, nil
}

func (sorter *defaultSorter) performParallelHorizontalSort(drawableImage draw.Image) error {
	yLength := drawableImage.Bounds().Dy()
	wg := sync.WaitGroup{}
	wg.Add(yLength)

	mu := sync.RWMutex{}
	iterationErrors := make(chan error, yLength)

	for y := 0; y < yLength; y += 1 {
		go func(yIndex int, errCh chan error) {
			defer wg.Done()

			mu.RLock()

			row, err := utils.GetImageRow(drawableImage, yIndex)
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
		}(y, iterationErrors)
	}

	wg.Wait()

	var err error = nil
	if len(iterationErrors) > 0 {
		err = <-iterationErrors
	}

	close(iterationErrors)
	return err
}

func (sorter *defaultSorter) performParallelVerticalSort(drawableImage draw.Image) error {
	xLength := drawableImage.Bounds().Dx()
	wg := sync.WaitGroup{}
	wg.Add(xLength)

	mu := sync.RWMutex{}
	iterationErrors := make(chan error, xLength)

	for x := 0; x < xLength; x += 1 {
		go func(xIndex int, errCh chan error) {
			defer wg.Done()

			mu.RLock()

			column, err := utils.GetImageColumn(drawableImage, xIndex)
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
		}(x, iterationErrors)
	}

	wg.Wait()

	var err error = nil
	if len(iterationErrors) > 0 {
		err = <-iterationErrors
	}

	close(iterationErrors)
	return err
}

// This is a helper function which performs the sorting of a given image strip by spliting it into intervals and sorting it by a given argument. This
// function is using the sorter preferences to determine how to perform the sort. There is also a maskedCoordinateFunc parameter which is a delegate
// used to retrieve information if a given pixel should be masked. We are using a external func for this in order to specify what coordinates should be
// looked up, beacuse this function has no access to the information which specific pixels from the image are processed now. Thanks to this approach, we
// can use a single function for both vertical and horizontal operations and just share a semi-fixed coordintes set.
func (sorter *defaultSorter) performSortOnImageStrip(imageStrip []color.Color, maskCoordinateFunc func(iteratedCoordinate int) (int, int)) ([]color.Color, error) {
	stripLength := len(imageStrip)
	sortedImageStrip := make([]color.Color, 0, stripLength)

	interval := CreateInterval(sorter.options.SortDeterminant)
	intervalMaxLength := sorter.calculateMaxIntervalLength()

	for x := 0; x < stripLength; x += 1 {
		currentColor := utils.ColorToRgba(imageStrip[x])

		isMasked, err := sorter.mask.IsMasked(maskCoordinateFunc(x))
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform a lookup to the mask image: %w", err)
		}

		// NOTE: isMasked and options dependecy solved using a quick K-Map
		passThrough := !isMasked || !sorter.options.UseMask

		if !utils.HasAnyTransparency(currentColor) && sorter.isMeetingIntervalRequirements(currentColor, isMasked, intervalMaxLength, interval) && passThrough {
			if err := interval.Append(currentColor); err != nil {
				return nil, fmt.Errorf("sorter: failed to append color to the interval: %w", err)
			}
		} else {
			if interval.Any() {
				sortedIntervalItems := interval.Sort(sorter.options.SortDirection, sorter.options.IntervalPainting)
				sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)

				interval = CreateInterval(sorter.options.SortDeterminant)
				intervalMaxLength = sorter.calculateMaxIntervalLength()
			}

			sortedImageStrip = append(sortedImageStrip, currentColor)
		}
	}

	if interval.Any() {
		sortedIntervalItems := interval.Sort(sorter.options.SortDirection, sorter.options.IntervalPainting)
		sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)
	}

	return sortedImageStrip, nil
}

func (sorter *defaultSorter) isMeetingIntervalRequirements(color color.RGBA, isMasked bool, maxLength int, interval Interval) bool {
	// NOTE: interval length and options dependecy solved using a quick K-Map
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

			h, _, _, _ := utils.ColorToHsla(color)
			hNorm := float64(h) / 360.0

			return hNorm >= lThreshold && hNorm <= uThreshold
		}
	case SplitBySaturation:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			_, s, _, _ := utils.ColorToHsla(color)
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

			abs := float64(int(color.R)*int(color.G)*int(color.B)) / 16581375.0

			return abs >= lThreshold && abs < uThreshold
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted interval determinant value")
	}
}

// This function determines the max interval length taking into account the randomness factor
func (sorter *defaultSorter) calculateMaxIntervalLength() int {
	if sorter.options.IntervalLength == 0 || sorter.options.IntervalLengthRandomFactor == 0 {
		return sorter.options.IntervalLength
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	factor := random.Intn(2*sorter.options.IntervalLengthRandomFactor) - sorter.options.IntervalLengthRandomFactor

	length := sorter.options.IntervalLength + factor
	if length < 1 {
		return 1
	} else {
		return length
	}
}
