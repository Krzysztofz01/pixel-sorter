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

// TODO: Implement support for the new mask adapter

type defaultSorter struct {
	image   image.Image
	mask    image.Image
	options *SorterOptions
}

func CreateSorter(image image.Image, mask image.Image, options *SorterOptions) (Sorter, error) {
	sorter := new(defaultSorter)
	sorter.image = image

	if mask != nil {
		drawableMask, err := utils.GetDrawableImage(mask)
		if err != nil {
			return nil, fmt.Errorf("sorter: can not get the drawable mask version: %w", err)
		}

		if ok, err := validateImageMask(sorter.image, drawableMask); !ok && err != nil {
			return nil, fmt.Errorf("sorter: validation of the mask failed: %w", err)
		}

		sorter.mask = drawableMask
	} else {
		sorter.mask = nil
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

		var rowMask []color.Color = nil
		if sorter.mask != nil {
			rowMask, err = utils.GetImageRow(sorter.mask, yIndex)
			if err != nil {
				return fmt.Errorf("sorter: failed to retrieve the image mask pixel row for a given index: %w", err)
			}
		}

		sortedRow, err := sorter.performSortOnImageStrip(row, rowMask)
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

			var rowMask []color.Color = nil
			if sorter.mask != nil {
				rowMask, err = utils.GetImageRow(sorter.mask, yIndex)
				if err != nil {
					errCh <- fmt.Errorf("sorter: failed to retrieve the image mask pixel row for a given index: %w", err)
					return
				}
			}

			sortedRow, err := sorter.performSortOnImageStrip(row, rowMask)
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

		var columnMask []color.Color = nil
		if sorter.mask != nil {
			columnMask, err = utils.GetImageColumn(sorter.mask, xIndex)
			if err != nil {
				return fmt.Errorf("sorter: failed to retrieve the image mask pixel column for a given index: %w", err)
			}
		}

		sortedColumn, err := sorter.performSortOnImageStrip(column, columnMask)
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

			var columnMask []color.Color = nil
			if sorter.mask != nil {
				columnMask, err = utils.GetImageColumn(sorter.mask, xIndex)
				if err != nil {
					errCh <- fmt.Errorf("sorter: failed to retrieve the image mask pixel column for a given index: %w", err)
					return
				}
			}

			sortedColumn, err := sorter.performSortOnImageStrip(column, columnMask)
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

// TODO: The isMaked code looks messy. We can extract the logic somehow in the future
// TODO: Instead of creating a mask strip we should implement a direct access to the mask
func (sorter *defaultSorter) performSortOnImageStrip(imageStrip []color.Color, maskStrip []color.Color) ([]color.Color, error) {
	stripLength := len(imageStrip)
	sortedImageStrip := make([]color.Color, 0, stripLength)
	sortDirection := GetSortDeterminantDirection(sorter.options.SortDeterminant)

	interval := sorter.CreateInterval()
	for x := 0; x < stripLength; x += 1 {
		currentColor, err := utils.ColorToRgba(imageStrip[x])
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to convert the given color to a RGBA struct representation: %w", err)
		}

		isMasked := false
		if maskStrip != nil {
			maskColor, err := utils.ColorToRgba(maskStrip[x])
			if err != nil {
				return nil, fmt.Errorf("sorter: failed to convert the given mask color to a RGBA struct representation: %w", err)
			}

			mR, mG, mB := utils.RgbaToIntComponents(maskColor)

			if mR == 255 && mG == 255 && mB == 255 {
				isMasked = true
			} else if mR == 0 && mG == 0 && mB == 0 {
				isMasked = false
			} else {
				return nil, errors.New("sorter: provied mask contains a invalid color")
			}
		}

		if !utils.HasAnyTransparency(currentColor) && sorter.isMeetingIntervalRequirements(currentColor) && !isMasked {
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

func (sorter *defaultSorter) isMeetingIntervalRequirements(color color.RGBA) bool {
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

func validateImageMask(img image.Image, mask image.Image) (bool, error) {
	iWidth := img.Bounds().Dx()
	mWidth := mask.Bounds().Dx()

	if iWidth != mWidth {
		return false, errors.New("sorter: the mask width is not matching the image width")
	}

	iHeight := img.Bounds().Dy()
	mHeight := mask.Bounds().Dy()

	if iHeight != mHeight {
		return false, errors.New("sorter: the mask height is not matching the image width")
	}

	// TODO: Check if mask colors are black or white. Adding support for the whole grayscale will make it possible to simplify this validation process
	for xIndex := 0; xIndex < mWidth; xIndex += 1 {
		for yIndex := 0; yIndex < mHeight; yIndex += 1 {
			color, err := utils.ColorToRgba(mask.At(xIndex, yIndex))
			if err != nil {
				return false, fmt.Errorf("sorter: failed to convert the color in the mask validation process: %w", err)
			}

			_, _, l := utils.RgbaToHsl(color)
			if l != 1.0 && l != 0.0 {
				return false, errors.New("sorter: the mask contains a invalid color")
			}
		}
	}

	return true, nil
}
