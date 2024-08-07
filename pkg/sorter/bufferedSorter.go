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

type bufferedSorter struct {
	image       *image.NRGBA
	maskImage   *image.NRGBA
	logger      SorterLogger
	cancel      func()
	cancelMutex sync.Mutex
	state       BufferedSorterState
}

// Create a new buffered image sorter instance by providing the image to be sorted and optional parameters such as mask image
// and a logger instance. This function will return a new buffered sorter instance or a error.
func CreateBufferedSorter(image image.Image, mask image.Image, logger SorterLogger) (BufferedSorter, error) {
	if image == nil {
		return nil, fmt.Errorf("sorter: can not create a sorter with the provided nil image")
	}

	sorter := new(bufferedSorter)
	sorter.cancel = nil
	sorter.cancelMutex = sync.Mutex{}
	sorter.state = CreateBufferedSorterState()
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

	return sorter, nil
}

func (sorter *bufferedSorter) CancelSort() bool {
	sorter.cancelMutex.Lock()
	defer sorter.cancelMutex.Unlock()

	if sorter.cancel == nil {
		return false
	}

	sorter.cancel()
	sorter.cancel = nil
	return true
}

func (sorter *bufferedSorter) Sort(options *SorterOptions) (image.Image, error) {
	var (
		srcImageNrgba       *image.NRGBA
		srcImageScaledNrgba *image.NRGBA
		srcMaskImageNrgba   *image.NRGBA
		srcImageRgba        *image.RGBA
		mask                Mask
		revertRotation      func(*image.NRGBA) *image.NRGBA
		sortingExecTime     time.Time = time.Now()
		err                 error     = nil
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer sorter.CancelSort()

	sorter.cancelMutex.Lock()
	sorter.cancel = cancel
	sorter.cancelMutex.Unlock()

	sorter.state.Apply(options)
	defer sorter.state.Rollback()

	if options.Scale != 1.0 {
		scalingExecTime := time.Now()

		if bufferedSrcImg, bufferedSrcMaskImg, ok := sorter.state.GetScaledImages(); ok {
			srcImageNrgba = bufferedSrcImg
			srcImageScaledNrgba = bufferedSrcImg
			srcMaskImageNrgba = bufferedSrcMaskImg
		} else {
			if srcImageNrgba, err = utils.ScaleImageNrgba(sorter.image, options.Scale); err != nil {
				return nil, fmt.Errorf("sorter: failed to scale the target image: %w", err)
			}

			if sorter.maskImage != nil {
				if srcMaskImageNrgba, err = utils.ScaleImageNrgba(sorter.maskImage, options.Scale); err != nil {
					return nil, fmt.Errorf("sorter: failed to scale the target image mask: %w", err)
				}
			} else {
				srcMaskImageNrgba = nil
			}

			srcImageScaledNrgba = srcImageNrgba
			sorter.state.SetScaledImages(srcImageNrgba, srcMaskImageNrgba)
		}

		sorter.logger.Debugf("Input images scaling took: %s", time.Since(scalingExecTime))
	} else {
		srcImageNrgba = sorter.image
		srcImageScaledNrgba = sorter.image
		srcMaskImageNrgba = sorter.maskImage
	}

	if options.Angle != 0 {
		if bufferedSrcImg, bufferedSrcMaskImg, ok := sorter.state.GetRotatedImages(); ok {
			revertRotationRectangle := srcImageNrgba.Rect
			srcImageNrgba = bufferedSrcImg
			srcMaskImageNrgba = bufferedSrcMaskImg

			revertRotation = func(n *image.NRGBA) *image.NRGBA {
				revertedNrgba := utils.RotateImageNrgba(n, -options.Angle)
				return utils.TrimImageTransparentWorkspaceNrgba(revertedNrgba, revertRotationRectangle)
			}
		} else {
			srcImageNrgba, revertRotation = utils.RotateImageWithRevertNrgba(srcImageNrgba, options.Angle)

			if srcMaskImageNrgba != nil {
				srcMaskImageNrgba = utils.RotateImageNrgba(srcImageNrgba, options.Angle)
			}

			sorter.state.SetRotatedImages(srcImageNrgba, srcMaskImageNrgba)
		}
	}

	if options.IntervalDeterminant == SplitByEdgeDetection {
		edgeDetectionExecTime := time.Now()

		if bufferedImage, ok := sorter.state.GetEdgeDetectionImage(); ok {
			srcMaskImageNrgba = bufferedImage
		} else {
			srcMaskImageNrgba, err := img.PerformEdgeDetection(srcImageNrgba, false, true)
			if err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the edge detection on the provided image: %w", err)
			}

			sorter.state.SetEdgeDetectionImage(srcMaskImageNrgba)
		}

		sorter.logger.Debugf("Edge detection took: %s.", time.Since(edgeDetectionExecTime))
	}

	if srcMaskImageNrgba != nil {
		maskExecTime := time.Now()
		if mask, err = CreateMaskFromNrgba(srcMaskImageNrgba); err != nil {
			return nil, fmt.Errorf("sorter: failed to create a new mask instance: %w", err)
		}

		sorter.logger.Debugf("Mask parsing took: %s.", time.Since(maskExecTime))
	} else {
		mask = CreateEmptyMask()
	}

	srcImageRgba = utils.NrgbaToRgbaImage(srcImageNrgba)
	dstImageRgba := utils.GetImageCopyRgba(srcImageRgba)

	for c := 0; c < options.Cycles; c += 1 {
		switch options.SortOrder {
		case SortVertical:
			{
				if err = performParallelColumnSorting(srcImageRgba, dstImageRgba, mask, options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical column sort: %w", err)
				}
			}
		case SortHorizontal:
			{
				if err = performParallelRowSorting(srcImageRgba, dstImageRgba, mask, options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal row sort: %w", err)
				}
			}
		case SortVerticalAndHorizontal:
			{
				if err = performParallelColumnSorting(srcImageRgba, dstImageRgba, mask, options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical column sort: %w", err)
				}

				copy(srcImageRgba.Pix, dstImageRgba.Pix)

				if err = performParallelRowSorting(srcImageRgba, dstImageRgba, mask, options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal row sort: %w", err)
				}
			}
		case SortHorizontalAndVertical:
			{
				if err = performParallelRowSorting(srcImageRgba, dstImageRgba, mask, options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the horizontal row sort: %w", err)
				}

				copy(srcImageRgba.Pix, dstImageRgba.Pix)

				if err = performParallelColumnSorting(srcImageRgba, dstImageRgba, mask, options, ctx); err != nil {
					return nil, fmt.Errorf("sorter: failed to perform the vertical column sort: %w", err)
				}
			}
		}

		if options.Cycles > 1 {
			copy(srcImageRgba.Pix, dstImageRgba.Pix)
		}
	}

	dstImageNrgba := utils.RgbaToNrgbaImage(dstImageRgba)
	if options.Angle != 0 {
		dstImageNrgba = revertRotation(dstImageNrgba)
	}

	switch options.Blending {
	case BlendingLighten:
		{
			if dstImageNrgba, err = utils.BlendImagesNrgba(srcImageScaledNrgba, dstImageNrgba, utils.LightenOnly); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the image lighten blending: %w", err)
			}
		}
	case BlendingDarken:
		{
			if dstImageNrgba, err = utils.BlendImagesNrgba(srcImageScaledNrgba, dstImageNrgba, utils.DarkenOnly); err != nil {
				return nil, fmt.Errorf("sorter: failed to perform the image darken blending: %w", err)
			}
		}
	case BlendingNone:
		break
	default:
		panic("sorter: invalid blending mode specified")
	}

	sorter.state.Commit()

	sorter.logger.Debugf("Pixel sorting took: %s.", time.Since(sortingExecTime))
	return dstImageNrgba, nil
}
