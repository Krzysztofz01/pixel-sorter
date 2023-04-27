package cmd

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	nestedFormatter "github.com/antonfisher/nested-logrus-formatter"
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
	FlagBlendingMode           string
	FlagVerboseLogging         bool
)

var (
	Logger      *logrus.Logger
	LocalLogger *logrus.Entry
)

var rootCmd = &cobra.Command{
	Use:   "pixel-sorter",
	Short: "Pixel sorting image editing utility implemented in Go.",
	Long:  "Pixel sorting image editing utility implemented in Go.",
}

func init() {
	Logger = CreateLogger()
	LocalLogger = CreateLocalLogger(Logger)

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

	rootCmd.PersistentFlags().StringVarP(&FlagBlendingMode, "blending-mode", "b", "none", "The blending mode algorithm to blend the sorted image into the original. Options: [none, lighten, darken].")
}

// Helper function used to validate and apply flag values into the sorter options struct
func parseCommonOptions() (*sorter.SorterOptions, error) {
	if FlagVerboseLogging {
		Logger.SetLevel(logrus.DebugLevel)
		Logger.SetReportCaller(true)

		LocalLogger = CreateLocalLogger(Logger)
	}

	options := sorter.GetDefaultSorterOptions()

	switch strings.ToLower(FlagSortDirection) {
	case "ascending":
		{
			options.SortDirection = sorter.SortAscending
		}
	case "descending":
		{
			options.SortDirection = sorter.SortDescending
		}
	case "random":
		{
			options.SortDirection = sorter.SortRandom
		}
	default:
		return nil, fmt.Errorf("invalid direction specified: %q", FlagSortDirection)
	}

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

	switch FlagBlendingMode {
	case "none":
		{
			options.Blending = sorter.BlendingNone
		}
	case "lighten":
		{
			options.Blending = sorter.BlendingLighten
		}
	case "darken":
		{
			options.Blending = sorter.BlendingDarken
		}
	default:
		{
			return nil, fmt.Errorf("invalid blending mode specified: %s", FlagBlendingMode)
		}
	}

	return options, nil
}

// Helper wrapper function used to perform the whole pixel sorting and IO operations according to the flags and provided options
func performPixelSorting(options *sorter.SorterOptions) error {
	if len(FlagImageFilePath) == 0 {
		return fmt.Errorf("invalid image path specified: %q", FlagImageFilePath)
	}

	format := strings.ToLower(FlagOutputFileType)
	if format != "jpg" && format != "png" {
		return fmt.Errorf("invalid output file format specified: %q", FlagOutputFileType)
	}

	img, err := utils.GetImageFromFile(FlagImageFilePath)
	if err != nil {
		return err
	}

	var mask image.Image = nil
	if len(FlagMaskFilePath) > 0 {
		mask, err = utils.GetImageFromFile(FlagMaskFilePath)
		if err != nil {
			return err
		}
	}

	sorter, err := sorter.CreateSorter(img, mask, Logger, options)
	if err != nil {
		return err
	}

	sortedImage, err := sorter.Sort()
	if err != nil {
		return err
	}

	if err := utils.StoreImageToFile(getOutputFileName(FlagImageFilePath), format, sortedImage); err != nil {
		return err
	}

	return nil
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
	LocalLogger.Info("Starting the pixel sorter.")
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		LocalLogger.Fatalf("Pixel sorting fatal failure: %s", err)
		os.Exit(1)
	}
	LocalLogger.Info("Pixel sorting finished.")
}

// Create a new instance of the logger
func CreateLogger() *logrus.Logger {
	formatter := &nestedFormatter.Formatter{
		TimestampFormat:  time.RFC3339Nano,
		HideKeys:         true,
		NoColors:         false,
		NoFieldsColors:   false,
		NoFieldsSpace:    false,
		ShowFullLevel:    false,
		NoUppercaseLevel: false,
		TrimMessages:     false,
		CallerFirst:      false,
	}

	return &logrus.Logger{
		Out:          os.Stdout,
		Formatter:    formatter,
		ReportCaller: false,
		Level:        logrus.InfoLevel,
	}
}

// Create a new prefixed logger entry instance
func CreateLocalLogger(logger *logrus.Logger) *logrus.Entry {
	return logger.WithField("prefix", "pixel-sorter-cli")
}
