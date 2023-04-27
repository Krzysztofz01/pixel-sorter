package cmd

import (
	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/spf13/cobra"
)

var saturationCmd = &cobra.Command{
	Use:   "saturation",
	Short: "Use saturation value as color sorting parameter.",
	Long:  "Use saturation value as color sorting parameter.",

	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := parseCommonOptions()
		if err != nil {
			LocalLogger.Errorf("Failed to parse the options: %s", err)
			return err
		}

		options.SortDeterminant = sorter.SortBySaturation
		if err := performPixelSorting(options); err != nil {
			LocalLogger.Errorf("Failed to perform the pixel sorting: %s", err)
		}

		return err
	},
}

func init() {
	hueCmd.SilenceUsage = true
	rootCmd.AddCommand(saturationCmd)
}
