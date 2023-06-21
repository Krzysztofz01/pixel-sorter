package cmd

import (
	"image"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Accept a video as the sorting input.",
	Long:  "Accept a video as the sorting input.",

	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := parseCommonOptions()
		if err != nil {
			LocalLogger.Errorf("Failed to parse the options: %s", err)
			return err
		}

		// TODO: Hardcoded for now. The sort determinant should be a flag and not a command
		options.SortDeterminant = sorter.SortByBrightness
		// TODO: Hardcoded for now. The explicit output path will be implemented in the future
		outputFileName := "sorter-video.mp4"

		var mask image.Image = nil
		if len(FlagMaskFilePath) > 0 {
			mask, err = utils.GetImageFromFile(FlagMaskFilePath)
			if err != nil {
				return err
			}
		}

		sorter, err := sorter.CreateVideoSorter(FlagImageFilePath, outputFileName, mask, Logger, options)
		if err != nil {
			return err
		}

		if err := sorter.Sort(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	videoCmd.SilenceUsage = true
	rootCmd.AddCommand(videoCmd)
}
