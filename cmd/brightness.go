package cmd

import (
	"fmt"
	"image"
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

		sorter, err := sorter.CreateSorter(img, mask, options)
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
	},
}

func init() {
	brightnessCmd.SilenceUsage = true
	rootCmd.AddCommand(brightnessCmd)
}
