package cmd

import (
	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/spf13/cobra"
)

var hueCmd = &cobra.Command{
	Use:   "hue",
	Short: "Use hue value as color sorting parameter.",
	Long:  "Use hue value as color sorting parameter.",

	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := parseCommonOptions()
		if err != nil {
			LocalLogger.Errorf("Failed to parse the options: %s", err)
			return err
		}

		options.SortDeterminant = sorter.SortByHue
		if err := performPixelSorting(options); err != nil {
			LocalLogger.Errorf("Failed to perform the pixel sorting: %s", err)
		}

		return err
	},
}

func init() {
	hueCmd.SilenceUsage = true
	rootCmd.AddCommand(hueCmd)
}
