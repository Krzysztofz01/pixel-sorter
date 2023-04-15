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
			return err
		}

		options.SortDeterminant = sorter.SortBySaturation
		return performPixelSorting(options)
	},
}

func init() {
	hueCmd.SilenceUsage = true
	rootCmd.AddCommand(saturationCmd)
}
