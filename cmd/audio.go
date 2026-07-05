package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var audioCmd = &cobra.Command{
	Use:   "audio",
	Short: "LucaSystem audio PAK tools",
	Long:  "Extract LucaSystem MUSIC/VOICE PAK audio and convert between native Ogg and MP3.",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var (
	AudioMP3       bool
	AudioDirection string
)

func init() {
	rootCmd.AddCommand(audioCmd)

	audioCmd.PersistentFlags().StringVarP(&PakInput, "input", "i", "", "input PAK/file/folder")
	audioCmd.PersistentFlags().StringVarP(&PakOutput, "output", "o", "", "output folder")
}

func printAudioSummary(prefix string, files, converted, skipped, errors int) {
	if converted > 0 || skipped > 0 {
		fmt.Printf("%s: %d native, %d converted, %d skipped, %d errors\n", prefix, files, converted, skipped, errors)
		return
	}
	fmt.Printf("%s: %d files, %d errors\n", prefix, files, errors)
}
