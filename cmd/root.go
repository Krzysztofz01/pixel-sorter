package cmd

import (
	"fmt"
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
	FlagInputMediaFilePath         string
	FlagOutputMediaFilePath        string
	FlagMaskImageFilePath          string
	FlagSortDeterminant            string
	FlagSortDirection              string
	FlagSortOrder                  string
	FlagIntervalDeterminant        string
	FlagIntervalLowerThreshold     float64
	FlagIntervalUpperThreshold     float64
	FlagAngle                      int
	FlagMask                       bool
	FlagIntervalLength             int
	FlagSortCycles                 int
	FlagImageScale                 float64
	FlagBlendingMode               string
	FlagVerboseLogging             bool
	FlagIntervalLengthRandomFactor int
)

var (
	Logger      *logrus.Logger
	LocalLogger *logrus.Entry
)

var Version string

var rootCmd = &cobra.Command{
	Use:   "pixel-sorter",
	Short: "Pixel sorting image editing utility implemented in Go.",
	Long:  fmt.Sprintf("Pixel sorting image editing utility implemented in Go.\nVersion: %s", Version),
}

func init() {
	Logger = CreateLogger()
	LocalLogger = CreateLocalLogger(Logger)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().BoolVarP(&FlagVerboseLogging, "verbose", "v", false, "Enable verbose logging mode.")

	rootCmd.PersistentFlags().StringVar(&FlagInputMediaFilePath, "input-media-path", "", "The path of the input media file to be processed.")
	rootCmd.MarkPersistentFlagRequired("input-media-path")

	rootCmd.PersistentFlags().StringVar(&FlagOutputMediaFilePath, "output-media-path", "", "The path of the output media file to be saved. The path should end with one of the supported extensions. [jpg, png]")
	rootCmd.MarkPersistentFlagRequired("output-media-path")

	rootCmd.PersistentFlags().StringVar(&FlagMaskImageFilePath, "mask-image-path", "", "The path of the mask image file used to process the input media.")

	rootCmd.PersistentFlags().StringVarP(&FlagSortDeterminant, "sort-determinant", "e", "brightness", "Parameter used as the argument for the sorting algorithm. Options: [brightness, hue, saturation].")

	rootCmd.PersistentFlags().StringVarP(&FlagSortDirection, "direction", "d", "ascending", "Pixel sorting direction in intervals. Options: [ascending, descending, shuffle, random].")

	rootCmd.PersistentFlags().StringVarP(&FlagSortOrder, "order", "o", "horizontal-vertical", "Order of the graphic sorting stages. Options: [horizontal, vertical, horizontal-vertical, vertical-horizontal].")

	rootCmd.PersistentFlags().StringVarP(&FlagIntervalDeterminant, "interval-determinant", "i", "brightness", "Parameter used to determine intervals. Options: [brightness, hue, saturation, mask, absolute, edge].")

	rootCmd.PersistentFlags().Float64VarP(&FlagIntervalLowerThreshold, "interval-lower-threshold", "l", 0.1, "The lower threshold of the interval determination process. Options: [0.0 - 1.0].")

	rootCmd.PersistentFlags().Float64VarP(&FlagIntervalUpperThreshold, "interval-upper-threshold", "u", 0.9, "The upper threshold of the interval determination process. Options: [0.0 - 1.0].")

	rootCmd.PersistentFlags().IntVarP(&FlagAngle, "angle", "a", 0, "The angle at which to sort the pixels.")

	rootCmd.PersistentFlags().BoolVarP(&FlagMask, "mask", "m", false, "Exclude the sorting effect from masked out ares of the image.")

	rootCmd.PersistentFlags().IntVarP(&FlagIntervalLength, "interval-max-length", "k", 0, "The max length of the interval. Zero means no length limits.")

	rootCmd.PersistentFlags().IntVarP(&FlagIntervalLengthRandomFactor, "interval-max-length-random-factor", "r", 0, "The value representing the range of values that can be randomly subtracted or added to the max interval length. Options: [>= 0]")

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

	if len(FlagInputMediaFilePath) == 0 {
		return nil, fmt.Errorf("cmd: invalid input media path specified (%s)", FlagInputMediaFilePath)
	}

	if len(FlagOutputMediaFilePath) == 0 {
		return nil, fmt.Errorf("cmd: invalid output media path specified (%s)", FlagOutputMediaFilePath)
	}

	options := sorter.GetDefaultSorterOptions()

	switch strings.ToLower(FlagSortDeterminant) {
	case "brightness":
		options.SortDeterminant = sorter.SortByBrightness
	case "hue":
		options.SortDeterminant = sorter.SortByHue
	case "saturation":
		options.SortDeterminant = sorter.SortBySaturation
	default:
		return nil, fmt.Errorf("cmd: invalid sort determinant specified (%s)", FlagSortDeterminant)
	}

	switch strings.ToLower(FlagSortDirection) {
	case "ascending":
		options.SortDirection = sorter.SortAscending
	case "descending":
		options.SortDirection = sorter.SortDescending
	case "shuffle":
		options.SortDirection = sorter.Shuffle
	case "random":
		options.SortDirection = sorter.SortRandom
	default:
		return nil, fmt.Errorf("cmd: invalid sort direction specified (%s)", FlagSortDirection)
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
		return nil, fmt.Errorf("cmd: invalid sort order specified (%s)", FlagSortOrder)
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
			if len(FlagMaskImageFilePath) == 0 {
				LocalLogger.Warnf("The interval determinant is using the mask, but not mask file has been specified.")
			}

			options.IntervalDeterminant = sorter.SplitByMask
		}
	case "absolute":
		options.IntervalDeterminant = sorter.SplitByAbsoluteColor
	case "edge":
		options.IntervalDeterminant = sorter.SplitByEdgeDetection
	default:
		return nil, fmt.Errorf("cmd: invalid interval determinant specified (%s)", FlagIntervalDeterminant)
	}

	switch FlagBlendingMode {
	case "none":
		options.Blending = sorter.BlendingNone
	case "lighten":
		options.Blending = sorter.BlendingLighten
	case "darken":
		options.Blending = sorter.BlendingDarken
	default:
		return nil, fmt.Errorf("cmd: invalid blending mode specified (%s)", FlagBlendingMode)
	}

	options.IntervalDeterminantUpperThreshold = FlagIntervalUpperThreshold
	options.IntervalDeterminantLowerThreshold = FlagIntervalLowerThreshold
	options.IntervalLength = FlagIntervalLength
	options.IntervalLengthRandomFactor = FlagIntervalLengthRandomFactor
	options.Angle = FlagAngle
	options.Cycles = FlagSortCycles
	options.Scale = FlagImageScale

	if FlagMask && len(FlagMaskImageFilePath) == 0 {
		LocalLogger.Warnf("The mask flag is set, but not mask file has been specified.")
	}

	options.UseMask = FlagMask

	if valid, msg := options.AreValid(); !valid {
		return nil, fmt.Errorf("cmd: %s", msg)
	}

	return options, nil
}

// Helper function used to determine if the current path file extension matches the possible extension collection.
func determineFileExtension(path string, extensions []string) (string, bool) {
	path, err := utils.EscapePathQuotes(path)
	if err != nil {
		return "", false
	}

	targetExtension := strings.ToLower(filepath.Ext(path))

	for _, allowedExtension := range extensions {
		lowerAllowedExtension := strings.ToLower(allowedExtension)

		if strings.HasSuffix(targetExtension, lowerAllowedExtension) {
			return lowerAllowedExtension, true
		}
	}

	return "", false
}

// Function used to execute the program (root command)
func Execute(args []string) {
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		LocalLogger.Fatalf("Pixel sorting fatal failure: %s", err)
		os.Exit(1)
	}
}

// TODO: Currently there is no support for setting the "image" command as the default one
func setDefaultCommand(args []string, defaultCommand string) []string {
	if len(args) > 1 {
		commands := make([]string, 5)
		for _, command := range rootCmd.Commands() {
			commands = append(commands, append(command.Aliases, command.Name())...)
		}

		commandSpecified := false
		currentCommand := args[0]
		for _, command := range commands {
			if command == currentCommand {
				commandSpecified = true
				break
			}
		}

		if !commandSpecified {
			args = append([]string{defaultCommand}, args...)
		}
	}

	return args
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
