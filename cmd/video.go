package cmd

import (
	"image"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Accept a video as the sorting input.",
	Long:  "Accept a video as the sorting input.",

	RunE: func(cmd *cobra.Command, args []string) error {
		LocalLogger.Info("Starting the video pixel sorting.")
		commandExecTime := time.Now()

		options, err := parseCommonOptions()
		if err != nil {
			LocalLogger.Errorf("Failed to parse the options: %s", err)
			return err
		}

		var mask image.Image = nil
		if len(FlagMaskImageFilePath) > 0 {
			mask, err = utils.GetImageFromFile(FlagMaskImageFilePath)
			if err != nil {
				return err
			}
		}

		sorter, err := sorter.CreateVideoSorter(FlagInputMediaFilePath, FlagOutputMediaFilePath, mask, SorterLogger, options)
		if err != nil {
			return err
		}

		if err := sorter.Sort(); err != nil {
			return err
		}

		LocalLogger.Infof("Video pixel sorting finished (%s).", time.Since(commandExecTime))
		return nil
	},
}

func init() {
	videoCmd.SilenceUsage = true
	rootCmd.AddCommand(videoCmd)
}
