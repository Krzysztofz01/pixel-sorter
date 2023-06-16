package sorter

import (
	"errors"
	"fmt"
	"image"

	vidio "github.com/AlexEidt/Vidio"
	"github.com/sirupsen/logrus"
)

type videoSorter struct {
	videoPath string
	maskImage image.Image
	logger    *logrus.Entry
	options   *SorterOptions
}

func CreateVideoSorter(videoPath string, mask image.Image, logger *logrus.Logger, options *SorterOptions) (Sorter, error) {
	// TODO: No options validation is made here. This logic should be moved from sorter factories to the SorterOptions struct.
	// TODO: The mask preparation and instance creation code should be shared and masks should be cached for video sorting.
	// TODO: The input path should be validated
	sorter := new(videoSorter)
	sorter.videoPath = videoPath
	sorter.maskImage = mask

	if logger == nil {
		return nil, errors.New("sorter: invalid logger reference provided")
	}

	sorter.logger = logger.WithField("prefix", "pixel-sorter")
	sorter.options = options

	return sorter, nil
}

func (sorter *videoSorter) Sort() (image.Image, error) {
	inputVideo, err := vidio.NewVideo(sorter.videoPath)
	if err != nil {
		return nil, fmt.Errorf("sorter: failed to access the input video file: %w", err)
	}

	videoOptions := vidio.Options{
		Bitrate: inputVideo.Bitrate(),
		FPS:     inputVideo.FPS(),
	}

	if inputVideo.HasStreams() {
		videoOptions.StreamFile = inputVideo.FileName()
	}

	// TODO: The output path is currently hardcoded
	outputVideo, err := vidio.NewVideoWriter("sorted.mp4", inputVideo.Width(), inputVideo.Height(), &videoOptions)
	if err != nil {
		return nil, fmt.Errorf("sorter: failed to access the output video file: %w", err)
	}

	defer outputVideo.Close()

	// TODO: Is it a better approach to use a uint8 buffer here?
	bufferImage := image.NewRGBA(image.Rect(0, 0, inputVideo.Width(), inputVideo.Height()))
	inputVideo.SetFrameBuffer(bufferImage.Pix)

	frameNumber := 1
	frameCount := inputVideo.Frames()
	for inputVideo.Read() {
		sorter.logger.Debugf("Frame %d/%d starting the processing step", frameNumber, frameCount)

		frameSorter, err := CreateSorter(bufferImage, sorter.maskImage, sorter.logger.Logger, sorter.options)
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to create a default sorter instance for frame %d/%d: %w", frameNumber, frameCount, err)
		}

		sortedFrame, err := frameSorter.Sort()
		if err != nil {
			return nil, fmt.Errorf("sorter: failed to perform the sorting for frame %d/%d: %w", frameNumber, frameCount, err)
		}

		sortedFrameRGBA, ok := sortedFrame.(*image.RGBA)
		if !ok {
			return nil, fmt.Errorf("sorter: failed to perform the frame buffer access cast for frame %d/%d: %w", frameNumber, frameCount, err)
		}

		if err := outputVideo.Write(sortedFrameRGBA.Pix); err != nil {
			return nil, fmt.Errorf("sorter: output write failure for frame %d/%d: %w", frameNumber, frameCount, err)
		}

		sorter.logger.Debugf("Frame %d/%d finishing the processing step", frameNumber, frameCount)
		frameNumber += 1
	}

	// TODO: This requries a different API (void/string return type?)
	return nil, nil
}
