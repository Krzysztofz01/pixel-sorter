package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/sirupsen/logrus"
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
	FlagSortCycles             int
	FlagImageScale             float64
	FlagVerboseLogging         bool
)

// TODO: Add verbose logging flag

var rootCmd = &cobra.Command{
	Use:   "pixel-sorter",
	Short: "Pixel sorting image editing utility implemented in Go.",
	Long:  "Pixel sorting image editing utility implemented in Go.",
}

func init() {
	logrus.SetFormatter(&customFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().BoolVarP(&FlagVerboseLogging, "verbose", "v", false, "Enable verbose logging mode.")

	rootCmd.PersistentFlags().StringVar(&FlagImageFilePath, "image-file-path", "", "The path of the image file to be processed.")
	rootCmd.MarkPersistentFlagRequired("image-file-path")

	rootCmd.PersistentFlags().StringVar(&FlagMaskFilePath, "mask-file-path", "", "The path of the image mask file to be process the image file.")

	rootCmd.PersistentFlags().StringVarP(&FlagOutputFileType, "output-format", "f", "jpg", "The output format of the graphic file. Options: [jpg, png].")

	rootCmd.PersistentFlags().StringVarP(&FlagSortDirection, "direction", "d", "ascending", "Pixel sorting direction in intervals. Options: [ascending, descending, random].")

	rootCmd.PersistentFlags().StringVarP(&FlagSortOrder, "order", "o", "horizontal-vertical", "Order of the graphic sorting stages. Options: [horizontal, vertical, horizontal-vertical, vertical-horizontal].")

	rootCmd.PersistentFlags().StringVarP(&FlagIntervalDeterminant, "interval-determinant", "i", "brightness", "Parameter used to determine intervals. Options: [brightness, hue, mask, absolute, edge].")
	rootCmd.PersistentFlags().Float64VarP(&FlagIntervalLowerThreshold, "interval-lower-threshold", "l", 0.1, "The lower threshold of the interval determination process. Options: [0.0 - 1.0].")
	rootCmd.PersistentFlags().Float64VarP(&FlagIntervalUpperThreshold, "interval-upper-threshold", "u", 0.9, "The upper threshold of the interval determination process. Options: [0.0 - 1.0].")

	rootCmd.PersistentFlags().IntVarP(&FlagAngle, "angle", "a", 0, "The angle at which to sort the pixels.")

	rootCmd.PersistentFlags().BoolVarP(&FlagMask, "mask", "m", false, "Exclude the sorting effect from masked out ares of the image.")

	rootCmd.PersistentFlags().IntVarP(&FlagIntervalLength, "interval-max-length", "k", 0, "The max length of the interval. Zero means no length limits.")

	rootCmd.PersistentFlags().IntVarP(&FlagSortCycles, "cycles", "c", 1, "The count of sorting cycles that should be performed on the image.")

	rootCmd.PersistentFlags().Float64VarP(&FlagImageScale, "scale", "s", 1, "Image downscaling percentage factor. Options: [0.0 - 1.0].")
}

// Helper function used to validate and apply flag values into the sorter options struct
func parseCommonOptions() (*sorter.SorterOptions, error) {
	if FlagVerboseLogging {
		logrus.SetLevel(logrus.DebugLevel)
	}

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
	case "saturation":
		options.IntervalDeterminant = sorter.SplitBySaturation
	case "mask":
		{
			if len(FlagMaskFilePath) == 0 {
				return nil, fmt.Errorf("invalid mask path specified")
			}
			options.IntervalDeterminant = sorter.SplitByMask
		}
	case "absolute":
		{
			options.IntervalDeterminant = sorter.SplitByAbsoluteColor
		}
	case "edge":
		{
			options.IntervalDeterminant = sorter.SplitByEdgeDetection
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

	if FlagSortCycles < 1 {
		return nil, fmt.Errorf("invalid cycles count specified")
	} else {
		options.Cycles = FlagSortCycles
	}

	if FlagImageScale < 0.0 || FlagImageScale > 1.0 {
		return nil, fmt.Errorf("invalid image scale percentage specified")
	} else {
		options.Scale = FlagImageScale
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

// Function used to execute the program (root command)
func Execute(args []string) {
	logrus.Info("Starting the pixel sorter.")
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Failure: %s", err)
		os.Exit(1)
	}
	logrus.Info("Pixel sorting finished.")
}

// Custom logrus formatter implementation
type customFormatter struct {
}

// Format func implementation for the custom logrus formatter
func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := "INF"
	switch entry.Level {
	case logrus.DebugLevel:
		level = "VER"
	case logrus.ErrorLevel:
		level = "ERR"
	case logrus.WarnLevel:
		level = "WRN"
	case logrus.InfoLevel:
		level = "INF"
	case logrus.FatalLevel:
		level = "ERR"
	}

	return []byte(fmt.Sprintf("[Pixel-Sorter] | [%s] | %s\n", level, entry.Message)), nil
}
