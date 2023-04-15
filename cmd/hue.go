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
			return err
		}

		options.SortDeterminant = sorter.SortByHue
		return performPixelSorting(options)
	},
}

func init() {
	hueCmd.SilenceUsage = true
	rootCmd.AddCommand(hueCmd)
}
