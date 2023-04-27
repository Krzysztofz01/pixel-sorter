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
			LocalLogger.Errorf("Failed to parse the options: %s", err)
			return err
		}

		options.SortDeterminant = sorter.SortByBrightness
		if err := performPixelSorting(options); err != nil {
			LocalLogger.Errorf("Failed to perform the pixel sorting: %s", err)
		}

		return err
	},
}

func init() {
	brightnessCmd.SilenceUsage = true
	rootCmd.AddCommand(brightnessCmd)
}
