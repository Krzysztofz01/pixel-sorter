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
func StoreImageToFile(fileName string, fileFormat string, img image.Image) error {
	if strings.ToLower(fileFormat) == "jpg" {
		file, err := os.Create(fmt.Sprintf("%s.jpg", fileName))
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

	} else if strings.ToLower(fileFormat) == "png" {
		file, err := os.Create(fmt.Sprintf("%s.png", fileName))
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
	} else {
		return errors.New("utils: specified file format is not supported")
	}
}
