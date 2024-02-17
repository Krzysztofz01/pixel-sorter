package sorter

import (
	"context"
	"errors"
	"fmt"
	"image"
	"sync"
	"time"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
)

type defaultVideoSorter struct {
	videoInputPath  string
	videoOutputPath string
	maskImage       image.Image
	logger          SorterLogger
	options         *SorterOptions
	cancel          func()
	cancelMutex     sync.Mutex
}

func CreateVideoSorter(inputVideoPath string, outputVideoPath string, mask image.Image, logger SorterLogger, options *SorterOptions) (VideoSorter, error) {
	sorter := new(defaultVideoSorter)
	// TODO: Currently there is not mask validation here. All the validation is happening for single frames.
	// Due to the fact thath the mask is stored here as a image and not as a parsed mask struct, it is very
	// inefficient. The mask parsing is called for all viedo frames.
	sorter.maskImage = mask
	sorter.cancel = nil
	sorter.cancelMutex = sync.Mutex{}

	if logger == nil {
		return nil, errors.New("sorter: invalid logger reference provided")
	}

	if logger == nil {
		sorter.logger = getDiscardLogger()
	} else {
		sorter.logger = logger
	}

	if options != nil {
		sorter.logger.Debugf("Running the video sorter with specified sorter options.")

		if valid, msg := options.AreValid(); !valid {
			sorter.logger.Debugf("Sorter options validation failed. Sorter options: %+v", *options)
			return nil, fmt.Errorf("sorter: %s", msg)
		}

		sorter.logger.Debugf("Sorter options validation passed. Sorter options: %+v", *options)
		sorter.options = options
	} else {
		sorter.logger.Debugf("Running the video sorter with default sorter options.")
		sorter.options = GetDefaultSorterOptions()
	}

	if len(inputVideoPath) == 0 {
		return nil, errors.New("sorter: invalid empty video input path specified")
	}

	if len(outputVideoPath) == 0 {
		return nil, errors.New("sorter: invalid empty video output path specified")
	}

	sorter.videoInputPath = inputVideoPath
	sorter.videoOutputPath = outputVideoPath

	return sorter, nil
}

func (sorter *defaultVideoSorter) Sort() error {
	sortingExecTime := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer sorter.CancelSort()

	sorter.cancelMutex.Lock()
	sorter.cancel = cancel
	sorter.cancelMutex.Unlock()

	inputVideo, err := vidio.NewVideo(sorter.videoInputPath)
	if err != nil {
		return fmt.Errorf("sorter: failed to access the input video file: %w", err)
	}

	videoOptions := &vidio.Options{
		Quality: 0,
		FPS:     inputVideo.FPS(),
	}

	if inputVideo.HasStreams() {
		videoOptions.StreamFile = inputVideo.FileName()
	}

	targetWidth := int(float64(inputVideo.Width()) * sorter.options.Scale)
	targetHeight := int(float64(inputVideo.Height()) * sorter.options.Scale)

	outputVideo, err := vidio.NewVideoWriter(sorter.videoOutputPath, targetWidth, targetHeight, videoOptions)
	if err != nil {
		return fmt.Errorf("sorter: failed to create the output video file: %w", err)
	}

	defer outputVideo.Close()

	frameBuffer := image.NewRGBA(image.Rect(0, 0, inputVideo.Width(), inputVideo.Height()))
	inputVideo.SetFrameBuffer(frameBuffer.Pix)

	frameNumber := 1
	frameCount := inputVideo.Frames()
	for inputVideo.Read() {
		select {
		case <-ctx.Done():
			return ErrSortingCancellation
		default:
		}

		sorter.logger.Debugf("Frame %d/%d starting the processing step", frameNumber, frameCount)
		frameSortingExecTime := time.Now()

		frameSorter, err := CreateSorter(frameBuffer, sorter.maskImage, nil, sorter.options)
		if err != nil {
			return fmt.Errorf("sorter: failed to create a default sorter instance for the frame %d/%d: %w", frameNumber, frameCount, err)
		}

		sortedFrame, err := frameSorter.Sort()
		if err != nil {
			return fmt.Errorf("sorter: failed to perform the sorting of the frame %d/%d: %w", frameNumber, frameCount, err)
		}

		sortedFrameRGBA := utils.ImageToRgbaImage(sortedFrame)

		if err := outputVideo.Write(sortedFrameRGBA.Pix); err != nil {
			return fmt.Errorf("sorter: failed to perform output video frame write for the frame %d/%d: %w", frameNumber, frameCount, err)
		}

		sorter.logger.Debugf("Frame %d/%d processing step finished in %s", frameNumber, frameCount, time.Since(frameSortingExecTime))
		frameNumber += 1
	}

	sorter.logger.Debugf("Video pixel sorting took: %s", time.Since(sortingExecTime))
	return nil
}

func (sorter *defaultVideoSorter) CancelSort() bool {
	sorter.cancelMutex.Lock()
	defer sorter.cancelMutex.Unlock()

	if sorter.cancel == nil {
		return false
	}

	sorter.cancel()
	sorter.cancel = nil
	return true
}