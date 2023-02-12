package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/spf13/cobra"
)

var (
	FlagImageFilePath          string
	FlagMaskFilePath           string
	FlagOutputFileType         string
	FlagSortDirection          string
	FlagSortOrder              string
	FlagIntervalDeterminant    string
	FlagIntervalLowerThreshold float64
	FlagIntervalUpperThreshold float64
	FlagAngle                  int
	FlagMask                   bool
	FlagIntervalLength         int
)

// TODO: Add verbose logging flag

var rootCmd = &cobra.Command{
	Use:   "pixel-sorter",
	Short: "Pixel sorting image editing utility implemented in Go.",
	Long:  "Pixel sorting image editing utility implemented in Go.",
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVar(&FlagImageFilePath, "image-file-path", "", "The path of the image file to be processed.")
	rootCmd.MarkPersistentFlagRequired("image-file-path")

	rootCmd.PersistentFlags().StringVar(&FlagMaskFilePath, "mask-file-path", "", "The path of the image mask file to be process the image file.")

	rootCmd.PersistentFlags().StringVarP(&FlagOutputFileType, "output-format", "f", "jpg", "The output format of the graphic file. Options: [jpg, png].")

	rootCmd.PersistentFlags().StringVarP(&FlagSortDirection, "direction", "d", "ascending", "Pixel sorting direction in intervals. Options: [ascending, descending, random].")

	rootCmd.PersistentFlags().StringVarP(&FlagSortOrder, "order", "o", "horizontal-vertical", "Order of the graphic sorting stages. Options: [horizontal, vertical, horizontal-vertical, vertical-horizontal].")

	rootCmd.PersistentFlags().StringVarP(&FlagIntervalDeterminant, "interval-determinant", "i", "brightness", "Parameter used to determine intervals. Options: [brightness, hue, mask].")
	rootCmd.PersistentFlags().Float64VarP(&FlagIntervalLowerThreshold, "interval-lower-threshold", "l", 0.1, "The lower threshold of the interval determination process. Options: [0.0 - 1.0].")
	rootCmd.PersistentFlags().Float64VarP(&FlagIntervalUpperThreshold, "interval-upper-threshold", "u", 0.9, "The upper threshold of the interval determination process. Options: [0.0 - 1.0].")

	rootCmd.PersistentFlags().IntVarP(&FlagAngle, "angle", "a", 0, "The angle at which to sort the pixels.")

	rootCmd.PersistentFlags().BoolVarP(&FlagMask, "mask", "m", false, "Exclude the sorting effect from masked out ares of the image.")

	rootCmd.PersistentFlags().IntVarP(&FlagIntervalLength, "interval-max-length", "k", 0, "The max length of the interval. Zero means no length limits.")
}

// Helper function used to validate and apply flag values into the sorter options struct
func parseCommonOptions() (*sorter.SorterOptions, error) {
	options := sorter.GetDefaultSorterOptions()

	switch strings.ToLower(FlagSortOrder) {
	case "horizontal":
		options.SortOrder = sorter.SortHorizontal
	case "vertical":
		options.SortOrder = sorter.SortVertical
	case "horizontal-vertical":
		options.SortOrder = sorter.SortHorizontalAndVertical
	case "vertical-horizontal":
		options.SortOrder = sorter.SortVerticalAndHorizontal
	default:
		return nil, fmt.Errorf("invalid order specified: %q", FlagSortOrder)
	}

	switch strings.ToLower(FlagIntervalDeterminant) {
	case "brightness":
		options.IntervalDeterminant = sorter.SplitByBrightness
	case "hue":
		options.IntervalDeterminant = sorter.SplitByHue
	case "mask":
		{
			if len(FlagMaskFilePath) == 0 {
				return nil, fmt.Errorf("invalid mask path specified")
			}
			options.IntervalDeterminant = sorter.SplitByMask
		}
	default:
		return nil, fmt.Errorf("invalid interval determinant specified: %q", FlagIntervalDeterminant)
	}

	if FlagIntervalLowerThreshold >= FlagIntervalUpperThreshold {
		return nil, fmt.Errorf("invalid interval thresholds: %f and: %f", FlagIntervalLowerThreshold, FlagIntervalUpperThreshold)
	}

	if FlagIntervalLowerThreshold < 0.0 || FlagIntervalLowerThreshold > 1.0 {
		return nil, fmt.Errorf("invalid lower interval threshold: %f", FlagIntervalLowerThreshold)
	} else {
		options.IntervalDeterminantLowerThreshold = FlagIntervalLowerThreshold
	}

	if FlagIntervalUpperThreshold < 0.0 || FlagIntervalUpperThreshold > 1.0 {
		return nil, fmt.Errorf("invalid lower interval threshold: %f", FlagIntervalUpperThreshold)
	} else {
		options.IntervalDeterminantUpperThreshold = FlagIntervalUpperThreshold
	}

	options.Angle = FlagAngle

	if FlagMask {
		if len(FlagMaskFilePath) == 0 {
			return nil, fmt.Errorf("invalid mask path specified")
		} else {
			options.UseMask = true
		}
	} else {
		options.UseMask = false
	}

	if FlagIntervalLength < 0 {
		return nil, fmt.Errorf("invalid interval length specified")
	} else {
		options.IntervalLength = FlagIntervalLength
	}

	return options, nil
}

// Helper function used to generate a output image file name based on the original file path
func getOutputFileName(inputFilePath string) string {
	fileName := filepath.Base(inputFilePath)
	fileNameParts := strings.Split(fileName, ".")
	if len(fileNameParts) != 2 {
		return "sorted"
	}

	return fmt.Sprintf("%s-sorted", fileNameParts[0])
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
