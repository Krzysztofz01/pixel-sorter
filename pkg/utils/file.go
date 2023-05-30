package utils

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

// Get the image from a file specified by the given path
func GetImageFromFile(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("utils: can not open the specified file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("utils: failed to decode the specified image: %w", err)
	}

	return img, nil
}

// Create a new file with the given name and format and store the given image in it
func StoreImageToFile(filePath string, fileFormat string, img image.Image) error {
	fileFormatLower := strings.ToLower(fileFormat)

	switch fileFormatLower {
	case "jpg", "jpeg":
		{
			file, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("utils: failed to create a new file: %w", err)
			}

			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()

			if err := jpeg.Encode(file, img, &jpeg.Options{Quality: 100}); err != nil {
				return fmt.Errorf("utils: failed to encode the image to jpeg: %w", err)
			}

			return nil
		}
	case "png":
		{
			file, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("utils: failed to create a new file: %w", err)
			}

			defer func() {
				if err := file.Close(); err != nil {
					panic(err)
				}
			}()

			if err := png.Encode(file, img); err != nil {
				return fmt.Errorf("utils: failed to encode the image to png: %w", err)
			}

			return nil
		}
	default:
		{
			return errors.New("utils: specified file format is not supported")
		}
	}
}
