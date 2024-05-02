package sorter

import (
	"image"
)

// TODO: Implement unit tests for the BufferedSorterState

type BufferedSorterState interface {
	// Apply sorter options changes to the buffered sorter state machine
	Apply(options *SorterOptions)

	// Confirm the incoming changes to the buffered sorter state machine
	Commit()

	// Reset non-commited incoming changes that came to the buffered sorter state machine
	Rollback()

	// Get the buffered scaled images set and a boolean value indicating if the values were buffered
	GetScaledImages() (*image.NRGBA, *image.NRGBA, bool)

	// Set the buffered scaled images associated to the incoming sorter options changes
	SetScaledImages(img, maskImage *image.NRGBA)

	// Get the buffered rotated images set and a boolean value indicating if the values were buffered
	GetRotatedImages() (*image.NRGBA, *image.NRGBA, bool)

	// Set the buffered rotated images associated to the incoming sorter options changes
	SetRotatedImages(img, maskImage *image.NRGBA)

	// Get the buffered edge detection image and a boolean value indicating if the value were buffered
	GetEdgeDetectionImage() (*image.NRGBA, bool)

	// Set the buffered edge detection image associated to the incoming sorter options changes
	SetEdgeDetectionImage(img *image.NRGBA)
}

func CreateBufferedSorterState() BufferedSorterState {
	return &bufferedSorterState{
		IncomingOptions:    nil,
		CurrentOptions:     nil,
		ImageScaled:        nil,
		ImageRotated:       nil,
		ImageEdgeDetection: nil,
		Commited:           false,
	}
}

type bufferedSorterState struct {
	IncomingOptions    *SorterOptions
	CurrentOptions     *SorterOptions
	ImageScaled        *BufferedPairEntry[*image.NRGBA]
	ImageRotated       *BufferedPairEntry[*image.NRGBA]
	ImageEdgeDetection *BufferedEntry[*image.NRGBA]
	Commited           bool
}

func (state *bufferedSorterState) Apply(options *SorterOptions) {
	if options == nil {
		panic("sorter: can not apply nil options as the incoming options for the buffered state")
	}

	if state.CurrentOptions == nil {
		state.CurrentOptions = options
	}

	state.IncomingOptions = options
	state.Commited = false
}

func (state *bufferedSorterState) Commit() {
	if state.IncomingOptions == nil {
		panic("sorter: can not apply nil options as the current options for the buffered state")
	}

	state.CurrentOptions = state.IncomingOptions
	state.IncomingOptions = nil
	state.Commited = true
}

func (state *bufferedSorterState) GetEdgeDetectionImage() (*image.NRGBA, bool) {
	if state.ImageEdgeDetection == nil {
		return nil, false
	}

	if state.CurrentOptions.Scale != state.IncomingOptions.Scale {
		return nil, false
	}

	if state.CurrentOptions.Angle != state.IncomingOptions.Angle {
		return nil, false
	}

	return state.ImageEdgeDetection.First, true
}

func (state *bufferedSorterState) GetRotatedImages() (*image.NRGBA, *image.NRGBA, bool) {
	if state.ImageRotated == nil {
		return nil, nil, false
	}

	if state.CurrentOptions.Scale != state.IncomingOptions.Scale {
		return nil, nil, false
	}

	if state.CurrentOptions.Angle != state.IncomingOptions.Angle {
		return nil, nil, false
	}

	return state.ImageRotated.First, state.ImageRotated.Second, true
}

func (state *bufferedSorterState) GetScaledImages() (*image.NRGBA, *image.NRGBA, bool) {
	if state.ImageScaled == nil || state.CurrentOptions.Scale != state.IncomingOptions.Scale {
		return nil, nil, false
	}

	return state.ImageScaled.First, state.ImageScaled.Second, true
}

func (state *bufferedSorterState) Rollback() {
	if state.Commited {
		return
	}

	state.IncomingOptions = nil
	state.CurrentOptions = nil
	state.ImageScaled = nil
	state.ImageRotated = nil
	state.ImageEdgeDetection = nil
	state.Commited = false
}

func (state *bufferedSorterState) SetEdgeDetectionImage(img *image.NRGBA) {
	state.ImageEdgeDetection = &BufferedEntry[*image.NRGBA]{
		First: img,
	}
}

func (state *bufferedSorterState) SetRotatedImages(img *image.NRGBA, maskImage *image.NRGBA) {
	state.ImageRotated = &BufferedPairEntry[*image.NRGBA]{
		First:  img,
		Second: maskImage,
	}
}

func (state *bufferedSorterState) SetScaledImages(img *image.NRGBA, maskImage *image.NRGBA) {
	state.ImageScaled = &BufferedPairEntry[*image.NRGBA]{
		First:  img,
		Second: maskImage,
	}
}

type BufferedEntry[TEntry any] struct {
	First TEntry
}

type BufferedPairEntry[TEntry any] struct {
	First  TEntry
	Second TEntry
}
