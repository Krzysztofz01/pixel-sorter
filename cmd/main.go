package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
)

var (
	argImagePath                         string
	argSortDeterminant                   string
	argIntervalDeterminant               string
	argIntervalDeterminantLowerThreshold float64
	argIntervalDeterminantUpperThreshold float64
	argAngle                             int
)

func main() {
	sorterOptions := sorter.GetDefaultSorterOptions()

	flag.StringVar(&argImagePath, "image-path", "", "Path to the image file to be sorted.")
	flag.StringVar(&argSortDeterminant, "sort-determinant", string(sorterOptions.SortDeterminant), "Pixel sorting algorithm parameter.")
	flag.StringVar(&argIntervalDeterminant, "interval-determinant", string(sorterOptions.IntervalDeterminant), "Sorting interval determination parameter.")
	flag.Float64Var(&argIntervalDeterminantLowerThreshold, "interval-lower-threshold", sorterOptions.IntervalDeterminantLowerThreshold, "The lower threshold for setting intervals.")
	flag.Float64Var(&argIntervalDeterminantUpperThreshold, "interval-upper-threshold", sorterOptions.IntervalDeterminantUpperThreshold, "The upper threshold for setting intervals.")
	flag.IntVar(&argAngle, "sorting-angle", sorterOptions.Angle, "Angle at which the pixels are to be sorted.")

	flag.Parse()

	if len(argImagePath) == 0 {
		fmt.Println("The path to the image must be specified.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// FIXME: Forced sorter options for development purposes
	sorterOptions.IntervalDeterminantLowerThreshold = 0.2
	sorterOptions.IntervalDeterminantUpperThreshold = 0.75
	// if err := applyArgumentsToSorterOptions(sorterOptions); err != nil {
	// 	fmt.Println("Invalid parameters provided.")
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }

	imageFile, err := os.Open(argImagePath)
	if err != nil {
		fmt.Println("Failed to open the image file.")
		fmt.Println(err)
		os.Exit(1)
	}

	// FIXME: Correct file handling
	defer imageFile.Close()

	image, _, err := image.Decode(imageFile)
	if err != nil {
		fmt.Println("Failed to decode the image file.")
		fmt.Println(err)
		os.Exit(1)
	}

	sorter, err := sorter.CreateSorter(image, sorterOptions)
	if err != nil {
		fmt.Println("Failed to initialize the image sorter utility.")
		fmt.Println(err)
		os.Exit(1)
	}

	sortedImage, err := sorter.Sort()
	if err != nil {
		fmt.Println("Failed to perform the image sorting process.")
		fmt.Println(err)
		os.Exit(1)
	}

	outputImageFile, err := os.Create(getOutputFileName(argImagePath))
	if err != nil {
		fmt.Println("Failed to create the output image file.")
		fmt.Println(err)
		os.Exit(1)
	}

	// FIXME: Correct file handling
	defer outputImageFile.Close()

	jpegOptions := jpeg.Options{Quality: 100}
	if err := jpeg.Encode(outputImageFile, sortedImage, &jpegOptions); err != nil {
		fmt.Println("Failed to encode the output image as JPEG format.")
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func applyArgumentsToSorterOptions(options *sorter.SorterOptions) error {
	return errors.New("cmd: not implemented")
}

func getOutputFileName(inputFilePath string) string {
	fileName := filepath.Base(inputFilePath)
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) != 2 {
		return "output.jpg"
	}

	return fmt.Sprintf("%s-sorted.jpg", fileNameParts[0])
}
