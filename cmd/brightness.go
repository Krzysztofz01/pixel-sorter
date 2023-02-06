package cmd

import (
	"fmt"
	"strings"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/spf13/cobra"
)

var brightnessCmd = &cobra.Command{
	Use:   "brightness",
	Short: "Use brightness value as color sorting parameter.",
	Long:  "Use brightness value as color sorting parameter.",

	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := parseCommonOptions()
		if err != nil {
			return err
		}

		switch strings.ToLower(FlagSortDirection) {
		case "ascending":
			options.SortDeterminant = sorter.SortByBrightnessAscending
		case "descending":
			options.SortDeterminant = sorter.SortByBrightnessDescending
		case "random":
			options.SortDeterminant = sorter.ShuffleByBrightness
		default:
			return fmt.Errorf("invalid direction specified: %q", FlagSortDirection)
		}

		if len(FlagFilePath) == 0 {
			return fmt.Errorf("invalid image path specified: %q", FlagFilePath)
		}

		format := strings.ToLower(FlagOutputFileType)
		if format != "jpg" && format != "png" {
			return fmt.Errorf("invalid output file format specified: %q", FlagOutputFileType)
		}

		image, err := utils.GetImageFromFile(FlagFilePath)
		if err != nil {
			return err
		}

		sorter, err := sorter.CreateSorter(image, options)
		if err != nil {
			return err
		}

		sortedImage, err := sorter.Sort()
		if err != nil {
			return err
		}

		if err := utils.StoreImageToFile(getOutputFileName(FlagFilePath), format, sortedImage); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(brightnessCmd)
}
