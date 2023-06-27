package cmd

import (
	"fmt"
	"image"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Perform a pixel sorting operation on the specified image file.",
	Long:  "Perform a pixel sorting operation on the specified image file.",

	RunE: func(cmd *cobra.Command, args []string) error {
		LocalLogger.Info("Starting the image pixel sorting.")
		commandExecTime := time.Now()

		options, err := parseCommonOptions()
		if err != nil {
			LocalLogger.Errorf("Failed to parse the options from the provided flags: %s", err)
			return err
		}

		format, ok := determineFileExtension(FlagOutputMediaFilePath, []string{"jpeg", "jpg", "png"})
		if !ok {
			return fmt.Errorf("cmd: invalid output image file format specified (%s)", FlagOutputMediaFilePath)
		}

		img, err := utils.GetImageFromFile(FlagInputMediaFilePath)
		if err != nil {
			return err
		}

		var mask image.Image = nil
		if len(FlagMaskImageFilePath) > 0 {
			mask, err = utils.GetImageFromFile(FlagMaskImageFilePath)
			if err != nil {
				return err
			}
		}

		sorter, err := sorter.CreateSorter(img, mask, Logger, options)
		if err != nil {
			return err
		}

		sortedImage, err := sorter.Sort()
		if err != nil {
			return err
		}

		if err := utils.StoreImageToFile(FlagOutputMediaFilePath, format, sortedImage); err != nil {
			return err
		}

		LocalLogger.Infof("Image pixel sorting finished (%s).", time.Since(commandExecTime))
		return nil
	},
}

func init() {
	imageCmd.SilenceUsage = true
	rootCmd.AddCommand(imageCmd)
}
