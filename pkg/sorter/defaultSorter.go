package sorter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"sync"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/img"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

type defaultSorter struct {
	image   image.Image
	mask    *Mask
	logger  SorterLogger
	options *SorterOptions
}

// Create a new image sorter instance by providing the image to be sorted and optional parameters such as mask image
// logger instance and custom sorter options. This function will return a new sorter instance or a error.
func CreateSorter(image image.Image, mask image.Image, logger SorterLogger, options *SorterOptions) (Sorter, error) {
	sorter := new(defaultSorter)
	sorter.image = image

	if logger == nil {
		sorter.logger = getDiscardLogger()
	} else {
		sorter.logger = logger
	}

	if options != nil {
		sorter.logger.Debugf("Running the sorter with specified sorter options.")

		if valid, msg := options.AreValid(); !valid {
			sorter.logger.Debugf("Sorter options validation failed. Sorter options: %+v", *options)
			return nil, fmt.Errorf("sorter: %s", msg)
		}

		sorter.logger.Debugf("Sorter options validation passed. Sorter options: %+v", *options)
		sorter.options = options
	} else {
		sorter.logger.Debugf("Running the sorter with default sorter options.")
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

		// NOTE: We are preventing the scaling when the IntervalDeterminant is set to edge detection, because the edge
		// detection based mask is already scaled, because it was performed on the scaled version of the  original image
		if sorter.options.Scale != 1.0 && sorter.options.IntervalDeterminant != SplitByEdgeDetection {
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

// Function used to iterate over image rows in parallel and invoking the image strip sorting on each row.
func (sorter *defaultSorter) performParallelRowSorting(src, dst *image.RGBA) error {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	wg := &sync.WaitGroup{}
	errt := utils.NewErrorTrap()

	for y := 0; y < height; y += 1 {
		wg.Add(1)
		go func(yIndex int) {
			defer wg.Done()

			if errt.IsSet() {
				return
			}

			if err := sorter.performImageStripSort(src, dst, 4*yIndex*height, 4*1, width); err != nil {
				errt.Set(fmt.Errorf("sorter: failed to perform image strip sorting for row %d: %w", yIndex, err))
				return
			}
		}(y)
	}

	wg.Wait()
	return errt.Err()
}

// Function used to iterate over image columns in paralle and ivoking the image strip sorting on each column
func (sorter *defaultSorter) performParallelColumnSorting(src, dst *image.RGBA) error {
	width := src.Bounds().Dx()
	height := src.Bounds().Dy()

	wg := &sync.WaitGroup{}
	errt := utils.NewErrorTrap()

	for x := 0; x < width; x += 1 {
		wg.Add(1)
		go func(xIndex int) {
			defer wg.Done()

			if errt.IsSet() {
				return
			}

			if err := sorter.performImageStripSort(src, dst, 4*xIndex, 4*width, height); err != nil {
				errt.Set(fmt.Errorf("sorter: failed to perform image strip sorting for column %d: %w", xIndex, err))
				return
			}
		}(x)
	}

	wg.Wait()
	return errt.Err()
}

// Deprecated: performParallelRowSorting
func (sorter *defaultSorter) performParallelHorizontalSort(drawableImage draw.Image) error {
	yLength := drawableImage.Bounds().Dy()
	wg := sync.WaitGroup{}
	wg.Add(yLength)

	mu := sync.RWMutex{}
	iterationErrors := make(chan error, yLength)

	for y := 0; y < yLength; y += 1 {
		/*go*/ func(yIndex int, errCh chan error) {
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

// Deprecated: performParallelColumnSorting
func (sorter *defaultSorter) performParallelVerticalSort(drawableImage draw.Image) error {
	xLength := drawableImage.Bounds().Dx()
	wg := sync.WaitGroup{}
	wg.Add(xLength)

	mu := sync.RWMutex{}
	iterationErrors := make(chan error, xLength)

	for x := 0; x < xLength; x += 1 {
		/*go*/ func(xIndex int, errCh chan error) {
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
// Deprecated: performImageStripSort
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
				// [ Changed interval API ] sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)
				for _, c := range sortedIntervalItems {
					sortedImageStrip = append(sortedImageStrip, c)
				}

				interval = CreateInterval(sorter.options.SortDeterminant)
				intervalMaxLength = sorter.calculateMaxIntervalLength()
			}

			sortedImageStrip = append(sortedImageStrip, currentColor)
		}
	}

	if interval.Any() {
		sortedIntervalItems := interval.Sort(sorter.options.SortDirection, sorter.options.IntervalPainting)
		// [ Changed interval API ] sortedImageStrip = append(sortedImageStrip, sortedIntervalItems...)
		for _, c := range sortedIntervalItems {
			sortedImageStrip = append(sortedImageStrip, c)
		}
	}

	return sortedImageStrip, nil
}

// Function used to sort a strip of pixels which can be a column or a row. The function accepts the source and destination image pointers. Due to the fact
// that the iteration is one-dimensional, we accept the start index and the iteration step size. The number of iteration steps is defined by the count. The
// function iterates over the strip and checks whether the interval requirements are met. If yes, they are appended to the interval, if not, they are
// written straight to the destination image. The intervals are also sorted and drawn into the image under some specific conditions.
func (sorter *defaultSorter) performImageStripSort(src, dst *image.RGBA, start, step, count int) error {
	var (
		buffer            []color.RGBA = make([]color.RGBA, 0, count)
		interval          Interval     = CreateInterval(sorter.options.SortDeterminant)
		intervalMaxLength int          = sorter.calculateMaxIntervalLength()
	)

	var (
		currentColor   color.RGBA
		isMasked       bool
		lowerThreshold float64
		upperThreshold float64
		err            error
	)

	for i, index := 0, start; i < count; i, index = i+1, index+step {
		currentColor.R = src.Pix[index+0]
		currentColor.G = src.Pix[index+1]
		currentColor.B = src.Pix[index+2]
		currentColor.A = src.Pix[index+3]

		// NOTE: Dont pass to interval if the pixel has any transparency
		if currentColor.A < 255 {
			goto sortAndResetInterval
		}

		// NOTE: Dont pass to interval if the interval max length has been reached (Solved using K-Map)
		if intervalMaxLength != 0 && interval.Count() >= intervalMaxLength {
			goto sortAndResetInterval
		}

		isMasked, err = sorter.mask.IsMaskedByIndex(index / 4)
		if err != nil {
			return fmt.Errorf("sorter: failed to perform a lookup to the mask image: %w", err)
		}

		// NOTE: Dont pass to interval if the mask is used and the pixel is masked (Solved using K-Map)
		if sorter.options.UseMask && isMasked {
			goto sortAndResetInterval
		}

		lowerThreshold = sorter.options.IntervalDeterminantLowerThreshold
		upperThreshold = sorter.options.IntervalDeterminantUpperThreshold

		// NOTE: Dont pass to interval if the interval determinant requirements are not meet
		if sorter.isMeetingIntervalDeterminant(currentColor, lowerThreshold, upperThreshold, isMasked) {
			goto sortAndResetInterval
		}

		if err := interval.Append(currentColor); err != nil {
			return fmt.Errorf("sorter: failed to append the current color to the interval: %w", err)
		}

		continue

	sortAndResetInterval:
		if interval.Any() {
			interval.SortToBuffer(sorter.options.SortDirection, sorter.options.IntervalPainting, &buffer)
			intervalMaxLength = sorter.calculateMaxIntervalLength()

			drawBufferIntoImage(dst, append(buffer, currentColor), index, step)
		} else {
			dst.Pix[index+0] = currentColor.R
			dst.Pix[index+1] = currentColor.G
			dst.Pix[index+2] = currentColor.B
			dst.Pix[index+3] = currentColor.A
		}
	}

	return nil
}

// Function used to check if the the given color is meeting the current interval determinant requirements taking the thresholds under account
func (sorter *defaultSorter) isMeetingIntervalDeterminant(c color.RGBA, lowerThreshold, upperThreshold float64, isMasked bool) bool {
	switch sorter.options.IntervalDeterminant {
	case SplitByBrightness:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			brightness := utils.CalculatePerceivedBrightness(c)
			return brightness >= lThreshold && brightness <= uThreshold
		}
	case SplitByHue:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			h, _, _, _ := utils.RgbaToHsla(c)
			hNorm := float64(h) / 360.0

			return hNorm >= lThreshold && hNorm <= uThreshold
		}
	case SplitBySaturation:
		{
			lThreshold := sorter.options.IntervalDeterminantLowerThreshold
			uThreshold := sorter.options.IntervalDeterminantUpperThreshold

			_, s, _, _ := utils.RgbaToHsla(c)
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

			abs := float64(int(c.R)*int(c.G)*int(c.B)) / 16581375.0

			return abs >= lThreshold && abs < uThreshold
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted interval determinant value")
	}
}

// Deprecated: isMeetingIntervalDeterminant
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

// Function used to calculate the max interval length by taking the options and randomness factor under account
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

// Function used to draw a color buffer to the destination image. The target position is determined by the iteration index and step value.
//
//lint:ignore U1000 Ignore unused false-positive caused by function call in label block
func drawBufferIntoImage(dst *image.RGBA, buffer []color.RGBA, index, step int) {
	var (
		color color.RGBA
	)

	for dstIndex, bufferIndex := index, len(buffer)-1; bufferIndex >= 0; dstIndex, bufferIndex = dstIndex-step, bufferIndex-1 {
		color = buffer[bufferIndex]

		dst.Pix[dstIndex+0] = color.R
		dst.Pix[dstIndex+1] = color.G
		dst.Pix[dstIndex+2] = color.B
		dst.Pix[dstIndex+3] = color.A
	}
}
