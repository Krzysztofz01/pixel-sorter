package sorter

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

// Function used to iterate over image rows in parallel and invoking the image strip sorting on each row.
func performParallelRowSorting(src, dst *image.RGBA, mask Mask, options *SorterOptions, ctx context.Context) error {
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

			if err := performImageStripSort(src, dst, mask, options, 4*yIndex*width, 4*1, width, ctx); err != nil {
				errt.Set(fmt.Errorf("sorter: failed to perform image strip sorting for row %d: %w", yIndex, err))
				ctx.Done()
				return
			}
		}(y)
	}

	wg.Wait()
	return errt.Err()
}

// Function used to iterate over image columns in paralle and ivoking the image strip sorting on each column
func performParallelColumnSorting(src, dst *image.RGBA, mask Mask, options *SorterOptions, ctx context.Context) error {
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

			if err := performImageStripSort(src, dst, mask, options, 4*xIndex, 4*width, height, ctx); err != nil {
				errt.Set(fmt.Errorf("sorter: failed to perform image strip sorting for column %d: %w", xIndex, err))
				ctx.Done()
				return
			}
		}(x)
	}

	wg.Wait()
	return errt.Err()
}

// Function used to sort a strip of pixels which can be a column or a row. The function accepts the source and destination image pointers. Due to the fact
// that the iteration is one-dimensional, we accept the start index and the iteration step size. The number of iteration steps is defined by the count. The
// function iterates over the strip and checks whether the interval requirements are met. If yes, they are appended to the interval, if not, they are
// written straight to the destination image. The intervals are also sorted and drawn into the image under some specific conditions.
func performImageStripSort(src, dst *image.RGBA, mask Mask, options *SorterOptions, start, step, count int, ctx context.Context) error {
	var (
		buffer                     []color.RGBA = make([]color.RGBA, 0, count)
		interval                   Interval     = CreateInterval(options.SortDeterminant)
		intervalLength             int          = options.IntervalLength
		intervalLengthRandomFactor int          = options.IntervalLengthRandomFactor
		lowerThreshold             float64      = options.IntervalDeterminantLowerThreshold
		upperThreshold             float64      = options.IntervalDeterminantUpperThreshold
		intervalMaxLength          int          = calculateMaxIntervalLength(intervalLength, intervalLengthRandomFactor)
	)

	var (
		currentColor color.RGBA
		isMasked     bool
		err          error
	)

	for i, index := 0, start; i < count; i, index = i+1, index+step {
		select {
		case <-ctx.Done():
			return ErrSortingCancellation
		default:
		}

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

		isMasked, err = mask.AtByIndexB(index / 4)
		if err != nil {
			return fmt.Errorf("sorter: failed to perform a lookup to the mask image: %w", err)
		}

		// NOTE: Dont pass to interval if the mask is used and the pixel is masked (Solved using K-Map)
		if options.UseMask && isMasked {
			goto sortAndResetInterval
		}

		// NOTE: Dont pass to interval if the interval determinant requirements are not meet
		if !isMeetingIntervalDeterminant(currentColor, options.IntervalDeterminant, lowerThreshold, upperThreshold, isMasked) {
			goto sortAndResetInterval
		}

		if err := interval.Append(currentColor); err != nil {
			return fmt.Errorf("sorter: failed to append the current color to the interval: %w", err)
		}

		continue

	sortAndResetInterval:
		if interval.Any() {
			buffer = buffer[:0]

			interval.SortToBuffer(options.SortDirection, options.IntervalPainting, &buffer)
			intervalMaxLength = calculateMaxIntervalLength(intervalLength, intervalLengthRandomFactor)

			drawBufferIntoImage(dst, append(buffer, currentColor), index, step)
		} else {
			dst.Pix[index+0] = currentColor.R
			dst.Pix[index+1] = currentColor.G
			dst.Pix[index+2] = currentColor.B
			dst.Pix[index+3] = currentColor.A
		}
	}

	// TODO: Incorporate this statement into the sortAndresetInterval label procedure
	if interval.Any() {
		buffer = buffer[:0]

		interval.SortToBuffer(options.SortDirection, options.IntervalPainting, &buffer)

		drawBufferIntoImage(dst, buffer, start+step*(count-1), step)
	}

	return nil
}

// Function used to check if the the given color is meeting the current interval determinant requirements taking the thresholds under account
func isMeetingIntervalDeterminant(c color.RGBA, determinant IntervalDeterminant, lowerThreshold, upperThreshold float64, isMasked bool) bool {
	switch determinant {
	case SplitByBrightness:
		{
			brightness := utils.CalculatePerceivedBrightness(c)

			return brightness >= lowerThreshold && brightness <= upperThreshold
		}
	case SplitByHue:
		{
			h, _, _, _ := utils.RgbaToHsla(c)
			hNorm := float64(h) / 360.0

			return hNorm >= lowerThreshold && hNorm <= upperThreshold
		}
	case SplitBySaturation:
		{
			_, s, _, _ := utils.RgbaToHsla(c)

			return s >= lowerThreshold && s <= upperThreshold
		}
	case SplitByMask, SplitByEdgeDetection:
		{
			return !isMasked
		}
	case SplitByAbsoluteColor:
		{
			abs := float64(int(c.R)*int(c.G)*int(c.B)) / 16581375.0

			return abs >= lowerThreshold && abs < upperThreshold
		}
	default:
		panic("sorter: invalid sorter state due to a corrupted interval determinant value")
	}
}

// Function used to calculate the max interval length by taking the options and randomness factor under account
func calculateMaxIntervalLength(intervalLength, intervalLengthRandomFactor int) int {
	if intervalLength == 0 || intervalLengthRandomFactor == 0 {
		return intervalLength
	}

	factor := utils.CIntn(2*intervalLengthRandomFactor) - intervalLengthRandomFactor

	length := intervalLength + factor
	if length < 1 {
		return 1
	} else {
		return length
	}
}

// Function used to draw a color buffer to the destination image. The target position is determined by the iteration
// index and step value. The specified index is the ending index.
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
