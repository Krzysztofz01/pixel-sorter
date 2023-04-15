package cmd

import (
	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
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

		options.SortDeterminant = sorter.SortByBrightness
		return performPixelSorting(options)
	},
}

func init() {
	brightnessCmd.SilenceUsage = true
	rootCmd.AddCommand(brightnessCmd)
}
