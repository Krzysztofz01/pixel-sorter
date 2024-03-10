package sorter

import (
	"context"
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/img"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

type defaultSorter struct {
	image       *image.NRGBA
	maskImage   *image.NRGBA
	mask        Mask
	logger      SorterLogger
	options     *SorterOptions
	cancel      func()
	cancelMutex sync.Mutex
}

// Create a new image sorter instance by providing the image to be sorted and optional parameters such as mask image
// logger instance and custom sorter options. This function will return a new sorter instance or a error.
func CreateSorter(image image.Image, mask image.Image, logger SorterLogger, options *SorterOptions) (Sorter, error) {
	if image == nil {
		return nil, fmt.Errorf("sorter: can not create a sorter with the provided nil image")
	}

	sorter := new(defaultSorter)
	sorter.cancel = nil
	sorter.cancelMutex = sync.Mutex{}
	sorter.image = utils.ImageToNrgbaImage(image)

	if mask != nil {
		sorter.maskImage = utils.ImageToNrgbaImage(mask)

		if image.Bounds() != mask.Bounds() {
			return nil, fmt.Errorf("sorter: can not create a sorter for a image and mask with bounds that are not matching")
		}
	}

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

	return sorter, nil
}

func (sorter *defaultSorter) CancelSort() bool {
	sorter.cancelMutex.Lock()
	defer sorter.cancelMutex.Unlock()

	if sorter.cancel == nil {
		return false
	}

	sorter.cancel()
	sorter.cancel = nil
	return true
}

func (sorter *defaultSorter) Sort() (image.Image, error) {
	var (
		srcImageNrgba   *image.NRGBA
		srcImageRgba    *image.RGBA
		maskImage       *image.NRGBA
		revertRotation  func(*image.NRGBA) *image.NRGBA
		sortingExecTime time.Time = time.Now()
		err             error     = nil
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer sorter.CancelSort()

	sorter.cancelMutex.Lock()
	sorter.cancel = cancel
	sorter.cancelMutex.Unlock()

	if sorter.options.Scale != 1.0 {
		scalingExecTime := time.Now()

		if srcImageNrgba, err = utils.ScaleImageNrgba(sorter.image, sorter.options.Scale); err != nil {
			return nil, fmt.Errorf("sorter: failed to scale the target image: %w", err)
		}

		if sorter.maskImage != nil {
			if maskImage, err = utils.ScaleImageNrgba(sorter.maskImage, sorter.options.Scale); err != nil {
				return nil, fmt.Errorf("sorter: failed to scale the target image mask: %w", err)
			}
		}

		sorter.logger.Debugf("Input images scaling took: %s", time.Since(scalingExecTime))
	} else {
		srcImageNrgba = sorter.image
		maskImage = sorter.maskImage
	}

	if sorter.options.Angle != 0 {
		srcImageNrgba, revertRotation = utils.RotateImageWithRevertNrgba(srcImageNrgba, sorter.options.Angle)

		if maskImage != nil {
			maskImage = utils.RotateImageNrgba(maskImage, sorter.options.Angle)
		}
	}

	if sorter.options.IntervalDeterminant == SplitByEdgeDetection {
		edgeDetectionExecTime := time.Now()
		maskImage, err = img.PerformEdgeDetection(srcImageNrgba, false, true)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform the edge detection on the provided image: %w", err)
		}

		sorter.logger.Debugf("Edge detection took: %s.", time.Since(edgeDetectionExecTime))
	}

	if maskImage != nil {
		maskExecTime := time.Now()
		if sorter.mask, err = CreateMaskFromNrgba(maskImage); err != nil {
			return nil, fmt.Errorf("sorter: failed to create a new mask instance: %w", err)
		}

		sorter.logger.Debugf("Mask parsing took: %s.", time.Since(maskExecTime))
	} else {
		sorter.mask = CreateEmptyMask()
	}

	srcImageRgba = utils.NrgbaToRgbaImage(srcImageNrgba)
	dstImageRgba := utils.GetImageCopyRgba(srcImageRgba)

	for c := 0; c < sorter.options.Cycles; c += 1 {
		switch sorter.options.SortOrder {
		case SortVertical:
			{
				if err = performParallelColumnSorting(srcImageRgba, dstImageRgba, sorter.mask, sorter.options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical column sort: %w", err)
				}
			}
		case SortHorizontal:
			{
				if err = performParallelRowSorting(srcImageRgba, dstImageRgba, sorter.mask, sorter.options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal row sort: %w", err)
				}
			}
		case SortVerticalAndHorizontal:
			{
				if err = performParallelColumnSorting(srcImageRgba, dstImageRgba, sorter.mask, sorter.options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical column sort: %w", err)
				}

				copy(srcImageRgba.Pix, dstImageRgba.Pix)

				if err = performParallelRowSorting(srcImageRgba, dstImageRgba, sorter.mask, sorter.options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal row sort: %w", err)
				}
			}
		case SortHorizontalAndVertical:
			{
				if err = performParallelRowSorting(srcImageRgba, dstImageRgba, sorter.mask, sorter.options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal row sort: %w", err)
				}

				copy(srcImageRgba.Pix, dstImageRgba.Pix)

				if err = performParallelColumnSorting(srcImageRgba, dstImageRgba, sorter.mask, sorter.options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical column sort: %w", err)
				}
			}
		}

		if sorter.options.Cycles > 1 {
			copy(srcImageRgba.Pix, dstImageRgba.Pix)
		}
	}

	dstImageNrgba := utils.RgbaToNrgbaImage(dstImageRgba)
	if sorter.options.Angle != 0 {
		dstImageNrgba = revertRotation(dstImageNrgba)
	}

	switch sorter.options.Blending {
	case BlendingLighten:
		{
			if dstImageNrgba, err = utils.BlendImagesNrgba(sorter.image, dstImageNrgba, utils.LightenOnly); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the image lighten blending: %w", err)
			}
		}
	case BlendingDarken:
		{
			if dstImageNrgba, err = utils.BlendImagesNrgba(sorter.image, dstImageNrgba, utils.DarkenOnly); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the image darken blending: %w", err)
			}
		}
	case BlendingNone:
		break
	default:
		panic("sorter: invalid blending mode specified")
	}

	sorter.logger.Debugf("Pixel sorting took: %s.", time.Since(sortingExecTime))
	return dstImageNrgba, nil
}
