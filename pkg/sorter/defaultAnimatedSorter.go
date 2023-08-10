package sorter

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"time"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/sirupsen/logrus"
)

const (
	animationFramesPerSeconds = 25
)

type defaultAnimatedSorter struct {
	videoOutputPath string
	image           image.Image
	mask            image.Image
	logger          *logrus.Entry
	keyframes       []*SorterOptions
}

// TODO: Better validation for input parameters
// TODO: Unit tests
func CreateAnimatedSorter(image image.Image, outputVideoPath string, mask image.Image, logger *logrus.Logger, options []*SorterOptions) (AnimatedSorter, error) {
	if image == nil {
		return nil, errors.New("sorter: invalid nil input image specified")
	}

	if len(outputVideoPath) == 0 {
		return nil, errors.New("sorter: invalid empty video output path specified")
	}

	if options == nil {
		return nil, errors.New("sorter: invalid nil sorter options collection specified")
	}

	if len(options) == 0 {
		return nil, errors.New("sorter: invalid empty sorter options collection specified")
	}

	if valid, msg := IsScaleValid(options); !valid {
		return nil, fmt.Errorf("sorter: invalid sorter scale due to %s", msg)
	}

	sorter := &defaultAnimatedSorter{
		videoOutputPath: outputVideoPath,
		image:           image,
		mask:            mask,
		keyframes:       options,
	}

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

	return sorter, nil
}

func (sorter *defaultAnimatedSorter) Sort() error {
	sortingExecTime := time.Now()

	videoOptions := &vidio.Options{
		Quality: 0,
		FPS:     animationFramesPerSeconds,
	}

	scale := sorter.GetScale()
	targetWidth := int(float64(sorter.image.Bounds().Dx()) * scale)
	targetHeight := int(float64(sorter.image.Bounds().Dy()) * scale)

	outputVideo, err := vidio.NewVideoWriter(sorter.videoOutputPath, targetWidth, targetHeight, videoOptions)
	if err != nil {
		return fmt.Errorf("sorter: failed to create the output video file: %w", err)
	}

	defer outputVideo.Close()

	for frameIndex, keyframe := range sorter.keyframes {
		sorter.logger.Debugf("Frame %d/%d starting the processing step", frameIndex+1, len(sorter.keyframes))
		frameSortingExecTime := time.Now()

		frameSorter, err := CreateSorter(sorter.image, sorter.mask, nil, keyframe)
		if err != nil {
			return fmt.Errorf("sorter: failed to create a default sorter instance for the frame %d/%d: %w", frameIndex+1, len(sorter.keyframes), err)
		}

		sortedFrame, err := frameSorter.Sort()
		if err != nil {
			return fmt.Errorf("sorter: failed to perform the sorting of the frame %d/%d: %w", frameIndex+1, len(sorter.keyframes), err)
		}

		sortedFrameRGBA, ok := sortedFrame.(*image.RGBA)
		if !ok {
			return fmt.Errorf("sorter: failed to perform the sorted frame buffer access cast for the frame %d/%d", frameIndex+1, len(sorter.keyframes))
		}

		if err := outputVideo.Write(sortedFrameRGBA.Pix); err != nil {
			return fmt.Errorf("sorter: failed to perform output video frame write for the frame %d/%d: %w", frameIndex+1, len(sorter.keyframes), err)
		}

		sorter.logger.Debugf("Frame %d/%d processing step finished in %s", frameIndex+1, len(sorter.keyframes), time.Since(frameSortingExecTime))
	}

	sorter.logger.Debugf("Animated pixel sorting took: %s", time.Since(sortingExecTime))
	return nil
}

func IsScaleValid(optionsCollection []*SorterOptions) (bool, string) {
	if optionsCollection == nil {
		return false, "options collection reference is nil"
	}

	if len(optionsCollection) == 0 {
		return false, "options collection is empty"
	}

	targetScale := optionsCollection[0].Scale
	for _, options := range optionsCollection[1:] {
		if targetScale != options.Scale {
			return false, "unconsistent scale value detected"
		}
	}

	return true, ""
}

func (sorter *defaultAnimatedSorter) GetScale() float64 {
	return sorter.keyframes[0].Scale
}
