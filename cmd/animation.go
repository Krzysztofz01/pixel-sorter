package cmd

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"os"
	"time"

	"github.com/Krzysztofz01/pixel-sorter/pkg/sorter"
	"github.com/Krzysztofz01/pixel-sorter/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	FlagKeyframesFilePath string
)

// FIXME: Broken keyframes flag valiation
// TODO: More consistent API for passing the keyframes or options
// TODO: Make the errors consistent accros all subcommands
var animationCmd = &cobra.Command{
	Use:   "animation",
	Short: "Accept a image as the sorting input and create a animation based on the specified options.",
	Long:  "Accept a image as the sorting input and create a animation based on the specified options.",

	RunE: func(cmd *cobra.Command, args []string) error {
		LocalLogger.Info("Starting the video pixel sorting.")
		commandExecTime := time.Now()

		keyframesFile, err := os.Open(FlagKeyframesFilePath)
		if err != nil {
			return fmt.Errorf("cmd: failed to open keyframes file: %w", err)
		}

		// TODO: Close() error check
		defer keyframesFile.Close()

		keyframesBytes, err := io.ReadAll(keyframesFile)
		if err != nil {
			return fmt.Errorf("cmd: failed to read th keyframes file content: %w", err)
		}

		keyframes := make([]*sorter.SorterOptions, 0)
		if err := json.Unmarshal(keyframesBytes, &keyframes); err != nil {
			return fmt.Errorf("cmd: failed to unmarshal the keyframes: %w", err)
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

		sorter, err := sorter.CreateAnimatedSorter(img, FlagOutputMediaFilePath, mask, Logger, keyframes)
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
	rootCmd.AddCommand(animationCmd)

	animationCmd.SilenceErrors = true

	animationCmd.Flags().StringVar(&FlagKeyframesFilePath, "keyframes-path", "", "The path of the keyframes file.")
	// TODO: The flag requirement is not working
	animationCmd.MarkFlagRequired("keyframes-path")

}
